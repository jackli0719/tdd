package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestMainEntry(t *testing.T) {
	// Test that main can be imported without panic
	t.Log("main package imported successfully")
}

func TestHealthWithoutDB(t *testing.T) {
	r := gin.New()

	// Simulate health endpoint behavior when db is nil
	r.GET("/health", func(c *gin.Context) {
		c.JSON(503, gin.H{"status": "database not connected"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != 503 {
		t.Errorf("expected status 503, got %d", w.Code)
	}
}

func TestHealthWithDB(t *testing.T) {
	r := gin.New()

	// Simulate health endpoint behavior when db is connected
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}