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

	err = db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{}, &model.OrderItem{}, &model.Staff{})
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
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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

func TestOrderService_GetByID_NotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

	_, err := svc.GetByID(999)
	if err != ErrOrderNotFound {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestOrderService_List(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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

	orders, total, err := svc.List(1, 10)
	if err != nil {
		t.Fatalf("failed to list orders: %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

func TestOrderService_Delete(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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

	err := svc.Delete(order.ID)
	if err != nil {
		t.Fatalf("failed to delete order: %v", err)
	}

	// Verify order is deleted
	_, err = svc.GetByID(order.ID)
	if err != ErrOrderNotFound {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestOrderService_Delete_InvalidState(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

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

	// Pay the order first
	svc.Paid(order.ID)

	// Cannot delete paid order
	err := svc.Delete(order.ID)
	if err != ErrInvalidOrderState {
		t.Errorf("expected ErrInvalidOrderState, got %v", err)
	}
}

func TestOrderService_Create_UserNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

	// Create product
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 100})

	_, err := svc.Create(&model.CreateOrderRequest{
		UserID: 999,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestOrderService_Create_ProductNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	_, err := svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 999, Quantity: 2},
		},
	})
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestOrderService_Create_InsufficientStock(t *testing.T) {
	db := setupOrderTestDB(t)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	svc := NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

	// Create user
	userRepo.Create(&model.User{Username: "testuser", Email: "test@example.com", Phone: "1234567890"})

	// Create product with low stock
	productRepo.Create(&model.Product{Name: "Test Product", Price: 99.99, Stock: 5})

	_, err := svc.Create(&model.CreateOrderRequest{
		UserID: 1,
		Items: []model.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 10},
		},
	})
	if err != ErrInsufficientStock {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}
