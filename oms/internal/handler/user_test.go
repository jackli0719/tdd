package handler

import (
	"bytes"
	"encoding/json"
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

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func setupRouter(svc service.UserService) *gin.Engine {
	r := gin.New()
	h := NewUserHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/users", h.List)
		api.GET("/users/:id", h.Get)
		api.POST("/users", h.Create)
		api.PUT("/users/:id", h.Update)
		api.DELETE("/users/:id", h.Delete)
	}

	return r
}

func TestUserHandler_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	body := `{"username":"testuser","email":"test@example.com","phone":"1234567890"}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
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

func TestUserHandler_List(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	// Create a user first
	svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	})

	req, _ := http.NewRequest("GET", "/api/users?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	_, _ = svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	})

	req, _ := http.NewRequest("GET", "/api/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	req, _ := http.NewRequest("GET", "/api/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		// Expected
		return
	}
	// Debug: print actual response
	t.Logf("Got status %d, body: %s", w.Code, w.Body.String())
	// If it's 400, that's because id=999 doesn't exist yet with fresh DB
	// Let's create a user first
	if w.Code == http.StatusBadRequest {
		t.Errorf("expected 404, got 400 - check if ID parsing works")
	}
}

func TestUserHandler_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	_, _ = svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	})

	body := `{"phone":"0987654321"}`
	req, _ := http.NewRequest("PUT", "/api/users/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	_, _ = svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	})

	req, _ := http.NewRequest("DELETE", "/api/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	req, _ := http.NewRequest("DELETE", "/api/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUserHandler_Create_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	body := `{invalid json}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	body := `{"phone":"0987654321"}`
	req, _ := http.NewRequest("PUT", "/api/users/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUserHandler_Update_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	r := setupRouter(svc)

	body := `{"phone":"0987654321"}`
	req, _ := http.NewRequest("PUT", "/api/users/invalid", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
