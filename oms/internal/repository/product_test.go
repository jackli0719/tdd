package repository

import (
	"fmt"

	"oms/internal/model"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupProductTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestProductRepository_Create(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewProductRepository(db)

	product := &model.Product{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}

	err := repo.Create(product)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if product.ID == 0 {
		t.Error("expected non-zero product ID")
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewProductRepository(db)

	product := &model.Product{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}
	repo.Create(product)

	found, err := repo.GetByID(product.ID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if found.Name != product.Name {
		t.Errorf("expected name %s, got %s", product.Name, found.Name)
	}
}

func TestProductRepository_GetByName(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewProductRepository(db)

	product := &model.Product{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}
	repo.Create(product)

	found, err := repo.GetByName("Test Product")
	if err != nil {
		t.Fatalf("failed to get product by name: %v", err)
	}

	if found.Price != product.Price {
		t.Errorf("expected price %f, got %f", product.Price, found.Price)
	}
}

func TestProductRepository_Update(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewProductRepository(db)

	product := &model.Product{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}
	repo.Create(product)

	product.Stock = 50
	err := repo.Update(product)
	if err != nil {
		t.Fatalf("failed to update product: %v", err)
	}

	found, _ := repo.GetByID(product.ID)
	if found.Stock != 50 {
		t.Errorf("expected stock 50, got %d", found.Stock)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewProductRepository(db)

	product := &model.Product{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	}
	repo.Create(product)

	err := repo.Delete(product.ID)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}

	_, err = repo.GetByID(product.ID)
	if err == nil {
		t.Error("expected error when getting deleted product")
	}
}

func TestProductRepository_List(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewProductRepository(db)

	for i := 0; i < 5; i++ {
		repo.Create(&model.Product{
			Name:  fmt.Sprintf("Test Product %d", i),
			Price: 99.99,
			Stock: 100,
		})
	}

	products, total, err := repo.List(0, 10, 0)
	if err != nil {
		t.Fatalf("failed to list products: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(products) != 5 {
		t.Errorf("expected 5 products, got %d", len(products))
	}
}
