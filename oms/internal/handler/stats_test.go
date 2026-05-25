package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"oms/internal/model"
	"oms/internal/repository"
	"oms/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupStatsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{}, &model.OrderItem{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func setupStatsRouter(svc service.OrderService) *gin.Engine {
	r := gin.New()
	h := NewStatsHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/stats/orders", h.OrderStats)
		api.GET("/stats/revenue", h.RevenueStats)
	}

	return r
}

func TestStatsHandler_OrderStats(t *testing.T) {
	db := setupStatsTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupStatsRouter(svc)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create products
	productRepo.Create(&model.Product{Name: "Product A", Price: 50.0, Stock: 100})
	productRepo.Create(&model.Product{Name: "Product B", Price: 30.0, Stock: 100})

	// Create orders
	svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items:  []model.CreateOrderItemRequest{{ProductID: 1, Quantity: 1}},
	})
	svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items:  []model.CreateOrderItemRequest{{ProductID: 2, Quantity: 2}},
	})

	req, _ := http.NewRequest("GET", "/api/stats/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStatsHandler_RevenueStats(t *testing.T) {
	db := setupStatsTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupStatsRouter(svc)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Product A", Price: 100.0, Stock: 100})

	// Create order
	svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items:  []model.CreateOrderItemRequest{{ProductID: 1, Quantity: 1}},
	})

	req, _ := http.NewRequest("GET", "/api/stats/revenue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStatsHandler_OrderStats_Empty(t *testing.T) {
	db := setupStatsTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupStatsRouter(svc)

	req, _ := http.NewRequest("GET", "/api/stats/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStatsHandler_RevenueStats_Empty(t *testing.T) {
	db := setupStatsTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupStatsRouter(svc)

	req, _ := http.NewRequest("GET", "/api/stats/revenue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}