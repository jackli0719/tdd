package service

import (
	"errors"
	"testing"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupReviewServiceTestDB(t *testing.T) (*gorm.DB, repository.ReviewRepository, repository.OrderRepository) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.Review{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db, repository.NewReviewRepository(db), repository.NewOrderRepository(db)
}

func TestReviewService_Create(t *testing.T) {
	db, reviewRepo, orderRepo := setupReviewServiceTestDB(t)
	svc := NewReviewService(reviewRepo, orderRepo)
	staffID := int64(3)
	order := &model.Order{
		OrderNo:     "ORD_REVIEW_1",
		UserID:      1,
		StaffID:     &staffID,
		TotalAmount: 100,
		Status:      model.OrderStatusCompleted,
	}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	review, err := svc.Create(1, &model.CreateReviewRequest{OrderID: order.ID, Rating: 5, Comment: "很好"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if review.UserID != 1 || review.StaffID != staffID {
		t.Fatalf("review = %+v, want user_id 1 staff_id %d", review, staffID)
	}

	_, err = svc.Create(1, &model.CreateReviewRequest{OrderID: order.ID, Rating: 4})
	if !errors.Is(err, ErrReviewAlreadyExists) {
		t.Fatalf("duplicate Create() error = %v, want ErrReviewAlreadyExists", err)
	}
}

func TestReviewService_CreateRejectsInvalidOrder(t *testing.T) {
	db, reviewRepo, orderRepo := setupReviewServiceTestDB(t)
	svc := NewReviewService(reviewRepo, orderRepo)
	staffID := int64(3)
	order := &model.Order{
		OrderNo:     "ORD_REVIEW_2",
		UserID:      1,
		StaffID:     &staffID,
		TotalAmount: 100,
		Status:      model.OrderStatusPaid,
	}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	_, err := svc.Create(1, &model.CreateReviewRequest{OrderID: order.ID, Rating: 5})
	if !errors.Is(err, ErrReviewOrderInvalid) {
		t.Fatalf("Create() error = %v, want ErrReviewOrderInvalid", err)
	}
}

func TestReviewService_CreateRejectsInvalidRating(t *testing.T) {
	_, reviewRepo, orderRepo := setupReviewServiceTestDB(t)
	svc := NewReviewService(reviewRepo, orderRepo)

	_, err := svc.Create(1, &model.CreateReviewRequest{OrderID: 1, Rating: 6})
	if !errors.Is(err, ErrInvalidRating) {
		t.Fatalf("Create() error = %v, want ErrInvalidRating", err)
	}
}
