package service

import (
	"errors"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrReviewNotFound      = errors.New("review not found")
	ErrReviewAlreadyExists = errors.New("review already exists for order")
	ErrReviewOrderInvalid  = errors.New("only completed orders with assigned staff can be reviewed")
	ErrInvalidRating       = errors.New("rating must be between 1 and 5")
)

// ReviewService handles review business logic.
type ReviewService struct {
	reviewRepo repository.ReviewRepository
	orderRepo  repository.OrderRepository
}

// NewReviewService creates a new ReviewService.
func NewReviewService(reviewRepo repository.ReviewRepository, orderRepo repository.OrderRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, orderRepo: orderRepo}
}

func (s *ReviewService) Create(userID int64, req *model.CreateReviewRequest) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, ErrInvalidRating
	}

	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	if order.Status != model.OrderStatusCompleted || order.StaffID == nil {
		return nil, ErrReviewOrderInvalid
	}

	if _, err := s.reviewRepo.GetByOrderID(req.OrderID); err == nil {
		return nil, ErrReviewAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	review := &model.Review{
		OrderID: req.OrderID,
		UserID:  userID,
		StaffID: *order.StaffID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetByID(id int64) (*model.Review, error) {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetByOrderID(orderID int64) (*model.Review, error) {
	review, err := s.reviewRepo.GetByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) ListByStaffID(staffID int64, page, pageSize int) ([]*model.Review, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	return s.reviewRepo.ListByStaffID(staffID, (page-1)*pageSize, pageSize)
}

func (s *ReviewService) List(page, pageSize int) ([]*model.Review, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	return s.reviewRepo.List((page-1)*pageSize, pageSize)
}

func (s *ReviewService) GetStaffSummary(staffID int64) (*model.StaffReviewSummary, error) {
	return s.reviewRepo.GetStaffSummary(staffID)
}

func (s *ReviewService) ListStaffSummaries() ([]model.StaffReviewSummary, error) {
	return s.reviewRepo.ListStaffSummaries()
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}
