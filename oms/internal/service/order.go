package service

import (
	"errors"
	"fmt"
	"time"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderState  = errors.New("invalid order state transition")
	ErrStaffNotAvailable  = errors.New("staff is not available")
	ErrAssignedStaffInvalid = errors.New("assigned staff not found or not available")
)

// OrderService handles order business logic
type OrderService interface {
	Create(req *model.CreateOrderRequest) (*model.Order, error)
	GetByID(id int64) (*model.Order, error)
	List(page, pageSize int) ([]*model.Order, int64, error)
	Delete(id int64) error
	AssignStaff(orderID int64, staffID *int64) error

	// State transitions
	Paid(id int64) error
	Ship(id int64) error
	Complete(id int64) error
	Cancel(id int64) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
	staffRepo   repository.StaffRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository,
	staffRepo repository.StaffRepository,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
		staffRepo:   staffRepo,
	}
}

func (s *orderService) Create(req *model.CreateOrderRequest) (*model.Order, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Verify staff if provided
	if req.StaffID != nil {
		staff, err := s.staffRepo.GetByID(*req.StaffID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAssignedStaffInvalid
			}
			return nil, err
		}
		if staff.Status != model.StaffStatusAvailable {
			return nil, ErrAssignedStaffInvalid
		}
	}

	// Calculate total and verify products
	var totalAmount float64
	items := make([]model.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrProductNotFound
			}
			return nil, err
		}

		if product.Stock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		items = append(items, model.OrderItem{
			ProductID: item.ProductID,
			Price:     product.Price,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	// Generate order number
	orderNo := fmt.Sprintf("ORD%d", time.Now().UnixNano())

	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      req.UserID,
		StaffID:     req.StaffID,
		TotalAmount: totalAmount,
		Status:      model.OrderStatusPending,
		Items:       items,
	}

	// Use transaction for order creation and stock decrement
	tx := s.orderRepo.DB().Begin()
	if err := s.orderRepo.CreateTx(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Decrement stock for each item
	for _, item := range req.Items {
		_, err := s.productRepo.DecrementStockTx(tx, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return nil, ErrInsufficientStock
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetByID(id int64) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
}

func (s *orderService) List(page, pageSize int) ([]*model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.orderRepo.List(offset, pageSize)
}

func (s *orderService) AssignStaff(orderID int64, staffID *int64) error {
	// Verify order exists
	_, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	// Verify staff if provided
	if staffID != nil {
		staff, err := s.staffRepo.GetByID(*staffID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAssignedStaffInvalid
			}
			return err
		}
		if staff.Status != model.StaffStatusAvailable {
			return ErrAssignedStaffInvalid
		}
	}

	return s.orderRepo.AssignStaff(orderID, staffID)
}

func (s *orderService) Delete(id int64) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	// Can only delete pending orders
	if order.Status != model.OrderStatusPending {
		return ErrInvalidOrderState
	}

	return s.orderRepo.Delete(id)
}

// State transitions
// pending -> paid -> shipped -> completed
//           -> cancelled (from pending or paid)

func (s *orderService) Paid(id int64) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status != model.OrderStatusPending {
		return ErrInvalidOrderState
	}

	order.Status = model.OrderStatusPaid
	return s.orderRepo.Update(order)
}

func (s *orderService) Ship(id int64) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status != model.OrderStatusPaid {
		return ErrInvalidOrderState
	}

	order.Status = model.OrderStatusShipped
	return s.orderRepo.Update(order)
}

func (s *orderService) Complete(id int64) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status != model.OrderStatusShipped {
		return ErrInvalidOrderState
	}

	order.Status = model.OrderStatusCompleted
	return s.orderRepo.Update(order)
}

func (s *orderService) Cancel(id int64) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusPaid {
		return ErrInvalidOrderState
	}

	order.Status = model.OrderStatusCancelled
	return s.orderRepo.Update(order)
}
