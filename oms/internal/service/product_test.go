package service

import (
	"fmt"

	"oms/internal/model"
	"oms/internal/repository"
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

	err = db.AutoMigrate(&model.Product{}, &model.Category{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func createTestCategory(t *testing.T, repo repository.CategoryRepository, name string) *model.Category {
	c := &model.Category{Name: name}
	repo.Create(c)
	found, _ := repo.GetByName(name)
	return found
}

func TestProductService_Create(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	req := &model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	}

	product, err := svc.Create(req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if product.Name != req.Name {
		t.Errorf("expected name %s, got %s", req.Name, product.Name)
	}
}

func TestProductService_Create_InvalidCategory(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	req := &model.CreateProductRequest{
		CategoryID: 999,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	}

	_, err := svc.Create(req)
	if err != ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestProductService_Create_DuplicateName(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	req := &model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	}

	svc.Create(req)

	_, err := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      49.99,
		Stock:      50,
	})

	if err != ErrProductExists {
		t.Errorf("expected ErrProductExists, got %v", err)
	}
}

func TestProductService_GetByID(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	})

	found, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if found.Name != created.Name {
		t.Errorf("expected name %s, got %s", created.Name, found.Name)
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	_, err := svc.GetByID(999)
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_Update(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	})

	stock := 50
	updated, err := svc.Update(created.ID, &model.UpdateProductRequest{
		Stock: &stock,
	})
	if err != nil {
		t.Fatalf("failed to update product: %v", err)
	}

	if updated.Stock != 50 {
		t.Errorf("expected stock 50, got %d", updated.Stock)
	}
}

func TestProductService_Update_ChangeCategory(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category1 := createTestCategory(t, categoryRepo, "Category 1")
	category2 := createTestCategory(t, categoryRepo, "Category 2")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category1.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	})

	newCategoryID := category2.ID
	updated, err := svc.Update(created.ID, &model.UpdateProductRequest{
		CategoryID: &newCategoryID,
	})
	if err != nil {
		t.Fatalf("failed to update product: %v", err)
	}

	if updated.CategoryID != category2.ID {
		t.Errorf("expected category_id %d, got %d", category2.ID, updated.CategoryID)
	}
}

func TestProductService_Update_InvalidCategory(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	})

	invalidCategoryID := int64(999)
	_, err := svc.Update(created.ID, &model.UpdateProductRequest{
		CategoryID: &invalidCategoryID,
	})
	if err != ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestProductService_Delete(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	})

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}

	_, err = svc.GetByID(created.ID)
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_List(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	for i := 0; i < 5; i++ {
		svc.Create(&model.CreateProductRequest{
			CategoryID: category.ID,
			Name:       fmt.Sprintf("Test Product %d", i),
			Price:      99.99,
			Stock:      100,
		})
	}

	products, total, err := svc.List(1, 10, 0)
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

func TestProductService_DecrementStock(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      100,
	})

	err := svc.DecrementStock(created.ID, 30)
	if err != nil {
		t.Fatalf("failed to decrement stock: %v", err)
	}

	product, _ := svc.GetByID(created.ID)
	if product.Stock != 70 {
		t.Errorf("expected stock 70, got %d", product.Stock)
	}
}

func TestProductService_DecrementStock_Insufficient(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Test Product",
		Price:      99.99,
		Stock:      10,
	})

	err := svc.DecrementStock(created.ID, 30)
	if err != ErrInsufficientStock {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}

func TestProductService_Update_DuplicateName(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	category := createTestCategory(t, categoryRepo, "Test Category")

	svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Product A",
		Price:      99.99,
		Stock:      100,
	})
	created, _ := svc.Create(&model.CreateProductRequest{
		CategoryID: category.ID,
		Name:       "Product B",
		Price:      49.99,
		Stock:      50,
	})

	_, err := svc.Update(created.ID, &model.UpdateProductRequest{
		Name: "Product A",
	})
	if err != ErrProductExists {
		t.Errorf("expected ErrProductExists, got %v", err)
	}
}

func TestProductService_Update_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	stock := 50
	_, err := svc.Update(999, &model.UpdateProductRequest{
		Stock: &stock,
	})
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_Delete_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	svc := NewProductService(productRepo, categoryRepo)

	err := svc.Delete(999)
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}
