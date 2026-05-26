package repository

import (
	"oms/internal/model"
	"time"

	"gorm.io/gorm"
)

// OrderRepository handles order data access
type OrderRepository interface {
	DB() *gorm.DB
	Create(order *model.Order) error
	CreateTx(tx *gorm.DB, order *model.Order) error
	GetByID(id int64) (*model.Order, error)
	GetByOrderNo(orderNo string) (*model.Order, error)
	Update(order *model.Order) error
	Delete(id int64) error
	List(offset, limit int) ([]*model.Order, int64, error)
	ListByUserID(userID int64, offset, limit int) ([]*model.Order, int64, error)
	ListByDateRange(start, end time.Time) ([]*model.Order, int64, error)
	AssignStaff(id int64, staffID *int64) error
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) DB() *gorm.DB {
	return r.db
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateTx(tx *gorm.DB, order *model.Order) error {
	return tx.Create(order).Error
}

func (r *orderRepository) GetByID(id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id int64) error {
	// Delete items first
	r.db.Where("order_id = ?", id).Delete(&model.OrderItem{})
	return r.db.Delete(&model.Order{}, id).Error
}

func (r *orderRepository) List(offset, limit int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	r.db.Model(&model.Order{}).Count(&total)
	err := r.db.Preload("Items").Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *orderRepository) ListByUserID(userID int64, offset, limit int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Preload("Items").Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *orderRepository) AssignStaff(id int64, staffID *int64) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("staff_id", staffID).Error
}

func (r *orderRepository) ListByDateRange(start, end time.Time) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	query := r.db.Model(&model.Order{}).Where("appointment_time >= ? AND appointment_time < ?", start, end)
	query.Count(&total)
	err := query.Where("appointment_time >= ? AND appointment_time < ?", start, end).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
