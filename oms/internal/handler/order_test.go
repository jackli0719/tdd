package handler

import (
	"bytes"
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

func setupOrderTestDB(t *testing.T) *gorm.DB {
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

func setupOrderRouter(svc service.OrderService) *gin.Engine {
	r := gin.New()
	h := NewOrderHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/orders", h.List)
		api.GET("/orders/:id", h.Get)
		api.POST("/orders", h.Create)
		api.DELETE("/orders/:id", h.Delete)
		api.POST("/orders/:id/paid", h.Paid)
		api.POST("/orders/:id/ship", h.Ship)
		api.POST("/orders/:id/complete", h.Complete)
		api.POST("/orders/:id/cancel", h.Cancel)
	}

	return r
}

func TestOrderHandler_Create(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupOrderRouter(svc)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	body := `{"user_id":1,"items":[{"product_id":1,"quantity":2}]}`
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOrderHandler_Get(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupOrderRouter(svc)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	// Create order
	svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})

	req, _ := http.NewRequest("GET", "/api/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOrderHandler_List(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupOrderRouter(svc)

	req, _ := http.NewRequest("GET", "/api/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOrderHandler_Paid(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupOrderRouter(svc)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	// Create order
	order, _ := svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})

	req, _ := http.NewRequest("POST", "/api/orders/1/paid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify status changed
	updatedOrder, _ := svc.GetByID(order.ID)
	if updatedOrder.Status != model.OrderStatusPaid {
		t.Errorf("expected status 'paid', got '%s'", updatedOrder.Status)
	}
}

func TestOrderHandler_Cancel(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(orderRepo, userRepo, productRepo)
	r := setupOrderRouter(svc)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	// Create order
	svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})

	req, _ := http.NewRequest("POST", "/api/orders/1/cancel", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
