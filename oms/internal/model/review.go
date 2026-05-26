package model

import "time"

// Review represents a user review for a completed order.
type Review struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int64     `json:"order_id" gorm:"not null;uniqueIndex:uk_reviews_order_id;index:idx_reviews_order_id"`
	UserID    int64     `json:"user_id" gorm:"not null;index:idx_reviews_user_id"`
	StaffID   int64     `json:"staff_id" gorm:"not null;index:idx_reviews_staff_id"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specifies the table name for Review.
func (Review) TableName() string {
	return "reviews"
}

// CreateReviewRequest is the request body for creating a review.
type CreateReviewRequest struct {
	OrderID int64  `json:"order_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,max=1000"`
}

// StaffReviewSummary contains aggregate review data for a staff member.
type StaffReviewSummary struct {
	StaffID      int64   `json:"staff_id"`
	AverageScore float64 `json:"average_score"`
	ReviewCount  int64   `json:"review_count"`
}
