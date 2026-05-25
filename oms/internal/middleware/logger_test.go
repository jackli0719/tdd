package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"oms/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	logger.Init("debug")
}

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("logs request with 200 status", func(t *testing.T) {
		router := gin.New()
		router.Use(Logger())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("logs request with 404 status", func(t *testing.T) {
		router := gin.New()
		router.Use(Logger())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Without a matching route, it should return 404
		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}