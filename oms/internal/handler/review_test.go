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

func setupReviewHandlerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.Review{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func setupReviewRouter(svc *service.ReviewService) *gin.Engine {
	r := gin.New()
	h := NewReviewHandler(svc)
	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	api := r.Group("/api")
	api.GET("/reviews", h.List)
	api.GET("/reviews/staff-summary", h.StaffSummary)
	api.GET("/reviews/:id", h.Get)
	api.POST("/reviews", h.Create)
	return r
}

func TestReviewHandler_Create(t *testing.T) {
	db := setupReviewHandlerTestDB(t)
	reviewRepo := repository.NewReviewRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewReviewService(reviewRepo, orderRepo)
	r := setupReviewRouter(svc)
	staffID := int64(2)
	order := &model.Order{OrderNo: "ORD_HANDLER_1", UserID: 1, StaffID: &staffID, TotalAmount: 10, Status: model.OrderStatusCompleted}
	db.Create(order)

	body := fmt.Sprintf(`{"order_id":%d,"rating":5,"comment":"很好"}`, order.ID)
	req, _ := http.NewRequest("POST", "/api/reviews", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	var resp response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("code = %d, want 0", resp.Code)
	}
}

func TestReviewHandler_CreateInvalidRating(t *testing.T) {
	db := setupReviewHandlerTestDB(t)
	reviewRepo := repository.NewReviewRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewReviewService(reviewRepo, orderRepo)
	r := setupReviewRouter(svc)

	req, _ := http.NewRequest("POST", "/api/reviews", bytes.NewBufferString(`{"order_id":1,"rating":6}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
