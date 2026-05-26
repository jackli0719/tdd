package repository

import (
	"strconv"
	"testing"

	"oms/internal/model"

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

func TestCategoryRepository_Create(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := NewCategoryRepository(db)

	category := &model.Category{
		Name:        "家政",
		Description: "家政服务类",
	}

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	if category.ID == 0 {
		t.Error("expected category ID to be set")
	}
}

func TestCategoryRepository_GetByID(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := NewCategoryRepository(db)

	category := &model.Category{
		Name:        "家政",
		Description: "家政服务类",
	}
	repo.Create(category)

	found, err := repo.GetByID(category.ID)
	if err != nil {
		t.Fatalf("failed to get category: %v", err)
	}

	if found.Name != "家政" {
		t.Errorf("expected name '家政', got '%s'", found.Name)
	}
}

func TestCategoryRepository_List(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := NewCategoryRepository(db)

	// Create 5 categories
	for i := 0; i < 5; i++ {
		repo.Create(&model.Category{
			Name:        "Category " + strconv.Itoa(i),
			Description: "Description",
		})
	}

	categories, total, err := repo.List(0, 10)
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

func TestCategoryRepository_HasProducts(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := NewCategoryRepository(db)

	// Create a category
	category := &model.Category{
		Name:        "家政",
		Description: "家政服务类",
	}
	repo.Create(category)

	// Check without products
	hasProducts, err := repo.HasProducts(category.ID)
	if err != nil {
		t.Fatalf("failed to check has products: %v", err)
	}
	if hasProducts {
		t.Error("expected no products, got true")
	}
}
