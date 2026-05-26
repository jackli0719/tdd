package repository

import (
	"oms/internal/model"

	"gorm.io/gorm"
)

// ReviewRepository handles review data access.
type ReviewRepository interface {
	Create(review *model.Review) error
	GetByID(id int64) (*model.Review, error)
	GetByOrderID(orderID int64) (*model.Review, error)
	ListByStaffID(staffID int64, offset, limit int) ([]*model.Review, int64, error)
	List(offset, limit int) ([]*model.Review, int64, error)
	GetStaffSummary(staffID int64) (*model.StaffReviewSummary, error)
	ListStaffSummaries() ([]model.StaffReviewSummary, error)
}

type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new ReviewRepository.
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByID(id int64) (*model.Review, error) {
	var review model.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByOrderID(orderID int64) (*model.Review, error) {
	var review model.Review
	if err := r.db.Where("order_id = ?", orderID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) ListByStaffID(staffID int64, offset, limit int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	query := r.db.Model(&model.Review{}).Where("staff_id = ?", staffID)
	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&reviews).Error; err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

func (r *reviewRepository) List(offset, limit int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	r.db.Model(&model.Review{}).Count(&total)
	if err := r.db.Offset(offset).Limit(limit).Order("id DESC").Find(&reviews).Error; err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

func (r *reviewRepository) GetStaffSummary(staffID int64) (*model.StaffReviewSummary, error) {
	summary := &model.StaffReviewSummary{StaffID: staffID}
	err := r.db.Model(&model.Review{}).
		Select("staff_id, COALESCE(AVG(rating), 0) AS average_score, COUNT(*) AS review_count").
		Where("staff_id = ?", staffID).
		Group("staff_id").
		Scan(summary).Error
	return summary, err
}

func (r *reviewRepository) ListStaffSummaries() ([]model.StaffReviewSummary, error) {
	var summaries []model.StaffReviewSummary
	err := r.db.Model(&model.Review{}).
		Select("staff_id, COALESCE(AVG(rating), 0) AS average_score, COUNT(*) AS review_count").
		Group("staff_id").
		Scan(&summaries).Error
	return summaries, err
}
