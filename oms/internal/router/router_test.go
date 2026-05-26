package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupRouterTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

func TestSetup_WithDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupRouterTestDB(t)
	r := gin.New()
	Setup(r, db, "test-secret")

	// Test health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSetup_WithoutDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	Setup(r, nil, "test-secret")

	// Test health endpoint without DB
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	// Test that API routes return 503 when no DB
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/users", nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w2.Code)
	}
}
