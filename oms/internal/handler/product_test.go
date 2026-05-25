package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"oms/internal/model"
	"oms/internal/repository"
	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
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

func setupProductRouter(svc service.ProductService) *gin.Engine {
	r := gin.New()
	h := NewProductHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/products", h.List)
		api.GET("/products/:id", h.Get)
		api.POST("/products", h.Create)
		api.PUT("/products/:id", h.Update)
		api.DELETE("/products/:id", h.Delete)
	}

	return r
}

func TestProductHandler_Create(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	body := `{"name":"Test Product","price":99.99,"stock":100}`
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
}

func TestProductHandler_List(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	// Create a product first
	svc.Create(&model.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	})

	req, _ := http.NewRequest("GET", "/api/products?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProductHandler_Get(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	svc.Create(&model.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	})

	req, _ := http.NewRequest("GET", "/api/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProductHandler_Get_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	req, _ := http.NewRequest("GET", "/api/products/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestProductHandler_Update(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	svc.Create(&model.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 100,
	})

	body := `{"stock":50}`
	req, _ := http.NewRequest("PUT", "/api/products/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProductHandler_Delete(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	created, _ := svc.Create(&model.CreateProductRequest{
		Name:  fmt.Sprintf("Test Product"),
		Price: 99.99,
		Stock: 100,
	})

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/products/%d", created.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProductHandler_Delete_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	req, _ := http.NewRequest("DELETE", "/api/products/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestProductHandler_Create_InvalidJSON(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	body := `{invalid json}`
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestProductHandler_Update_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	body := `{"stock":50}`
	req, _ := http.NewRequest("PUT", "/api/products/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestProductHandler_InvalidID(t *testing.T) {
	db := setupProductTestDB(t)
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	r := setupProductRouter(svc)

	req, _ := http.NewRequest("GET", "/api/products/invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
