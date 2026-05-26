package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"oms/internal/model"
	"oms/internal/repository"
	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
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

func setupCategoryRouter(svc service.CategoryService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewCategoryHandler(svc)
	r.GET("/api/categories", handler.List)
	r.GET("/api/categories/:id", handler.Get)
	r.POST("/api/categories", handler.Create)
	r.PUT("/api/categories/:id", handler.Update)
	r.DELETE("/api/categories/:id", handler.Delete)

	return r
}

func TestCategoryHandler_Create(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	r := setupCategoryRouter(svc)

	body := `{"name":"家政","description":"家政服务类"}`
	req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBufferString(body))
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

func TestCategoryHandler_List(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	r := setupCategoryRouter(svc)

	// Create a category first
	svc.Create(&model.CreateCategoryRequest{
		Name:        "家政",
		Description: "家政服务类",
	})

	req, _ := http.NewRequest("GET", "/api/categories?page=1&page_size=10", nil)
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

func TestCategoryHandler_Get(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	r := setupCategoryRouter(svc)

	// Create a category first
	category, _ := svc.Create(&model.CreateCategoryRequest{
		Name:        "家政",
		Description: "家政服务类",
	})

	req, _ := http.NewRequest("GET", "/api/categories/"+strconv.FormatInt(category.ID, 10), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryHandler_Delete(t *testing.T) {
	db := setupCategoryTestDB(t)
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	r := setupCategoryRouter(svc)

	// Create a category first
	category, _ := svc.Create(&model.CreateCategoryRequest{
		Name:        "家政",
		Description: "家政服务类",
	})

	req, _ := http.NewRequest("DELETE", "/api/categories/"+strconv.FormatInt(category.ID, 10), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
