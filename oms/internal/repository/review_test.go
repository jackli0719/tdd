package repository

import (
	"testing"

	"oms/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupReviewRepositoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Review{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestReviewRepository_CreateAndGetByOrderID(t *testing.T) {
	db := setupReviewRepositoryTestDB(t)
	repo := NewReviewRepository(db)

	review := &model.Review{OrderID: 1, UserID: 2, StaffID: 3, Rating: 5, Comment: "很好"}
	if err := repo.Create(review); err != nil {
		t.Fatalf("failed to create review: %v", err)
	}
	if review.ID == 0 {
		t.Fatal("expected non-zero review ID")
	}

	found, err := repo.GetByOrderID(1)
	if err != nil {
		t.Fatalf("failed to get review by order id: %v", err)
	}
	if found.Rating != 5 {
		t.Fatalf("rating = %d, want 5", found.Rating)
	}
}

func TestReviewRepository_ListByStaffIDAndSummary(t *testing.T) {
	db := setupReviewRepositoryTestDB(t)
	repo := NewReviewRepository(db)

	repo.Create(&model.Review{OrderID: 1, UserID: 1, StaffID: 7, Rating: 4})
	repo.Create(&model.Review{OrderID: 2, UserID: 2, StaffID: 7, Rating: 5})
	repo.Create(&model.Review{OrderID: 3, UserID: 3, StaffID: 8, Rating: 2})

	reviews, total, err := repo.ListByStaffID(7, 0, 10)
	if err != nil {
		t.Fatalf("failed to list reviews: %v", err)
	}
	if len(reviews) != 2 || total != 2 {
		t.Fatalf("got len=%d total=%d, want 2/2", len(reviews), total)
	}

	summary, err := repo.GetStaffSummary(7)
	if err != nil {
		t.Fatalf("failed to get staff summary: %v", err)
	}
	if summary.ReviewCount != 2 || summary.AverageScore != 4.5 {
		t.Fatalf("summary = %+v, want count 2 average 4.5", summary)
	}
}
