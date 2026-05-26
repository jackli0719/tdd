package service

import (
	"strconv"
	"testing"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCategoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&model.Category{}, &model.Product{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestCategoryService_Create(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := NewCategoryService(repo)

	category, err := svc.Create(&model.CreateCategoryRequest{
		Name:        "家政",
		Description: "家政服务类",
	})
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	if category.Name != "家政" {
		t.Errorf("expected name '家政', got '%s'", category.Name)
	}
}

func TestCategoryService_List(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := NewCategoryService(repo)

	// Create 5 categories
	for i := 0; i < 5; i++ {
		svc.Create(&model.CreateCategoryRequest{
			Name:        "Category " + strconv.Itoa(i),
			Description: "Description",
		})
	}

	categories, total, err := svc.List(1, 10)
	if err != nil {
		t.Fatalf("failed to list categories: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(categories) != 5 {
		t.Errorf("expected 5 categories, got %d", len(categories))
	}
}

func TestCategoryService_Delete_WithProducts(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := NewCategoryService(repo)

	// Create a category
	category, _ := svc.Create(&model.CreateCategoryRequest{
		Name:        "家政",
		Description: "家政服务类",
	})

	// Create a product in this category
	productRepo := repository.NewProductRepository(db)
	productRepo.Create(&model.Product{
		CategoryID: category.ID,
		Name:       "4小时保洁",
		Price:      98,
		Stock:      10,
	})

	// Try to delete category - should fail
	err := svc.Delete(category.ID)
	if err != ErrCategoryHasProducts {
		t.Errorf("expected ErrCategoryHasProducts, got %v", err)
	}
}
