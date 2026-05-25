package service

import (
	"oms/internal/model"
	"oms/internal/repository"
	"testing"

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

func TestOrderService_StateTransitions(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	// Create order
	order, err := svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	if order.Status != model.OrderStatusPending {
		t.Errorf("expected status 'pending', got '%s'", order.Status)
	}

	// Test Paid transition
	err = svc.Paid(order.ID)
	if err != nil {
		t.Fatalf("failed to pay order: %v", err)
	}

	order, _ = svc.GetByID(order.ID)
	if order.Status != model.OrderStatusPaid {
		t.Errorf("expected status 'paid', got '%s'", order.Status)
	}

	// Test Ship transition
	err = svc.Ship(order.ID)
	if err != nil {
		t.Fatalf("failed to ship order: %v", err)
	}

	order, _ = svc.GetByID(order.ID)
	if order.Status != model.OrderStatusShipped {
		t.Errorf("expected status 'shipped', got '%s'", order.Status)
	}

	// Test Complete transition
	err = svc.Complete(order.ID)
	if err != nil {
		t.Fatalf("failed to complete order: %v", err)
	}

	order, _ = svc.GetByID(order.ID)
	if order.Status != model.OrderStatusCompleted {
		t.Errorf("expected status 'completed', got '%s'", order.Status)
	}
}

func TestOrderService_CancelFromPending(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo)

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

	// Cancel from pending
	err := svc.Cancel(order.ID)
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}

	order, _ = svc.GetByID(order.ID)
	if order.Status != model.OrderStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", order.Status)
	}
}

func TestOrderService_CancelFromPaid(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	// Create order and pay
	order, _ := svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})
	svc.Paid(order.ID)

	// Cancel from paid
	err := svc.Cancel(order.ID)
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}

	order, _ = svc.GetByID(order.ID)
	if order.Status != model.OrderStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", order.Status)
	}
}

func TestOrderService_InvalidTransitions(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo)

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

	// Cannot ship from pending
	err := svc.Ship(order.ID)
	if err != ErrInvalidOrderState {
		t.Errorf("expected ErrInvalidOrderState, got %v", err)
	}

	// Cannot complete from pending
	err = svc.Complete(order.ID)
	if err != ErrInvalidOrderState {
		t.Errorf("expected ErrInvalidOrderState, got %v", err)
	}

	// Cannot cancel from shipped
	svc.Paid(order.ID)
	svc.Ship(order.ID)
	err = svc.Cancel(order.ID)
	if err != ErrInvalidOrderState {
		t.Errorf("expected ErrInvalidOrderState, got %v", err)
	}
}
