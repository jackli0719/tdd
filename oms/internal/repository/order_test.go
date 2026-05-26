package repository

import (
	"fmt"

	"oms/internal/model"
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

	err = db.AutoMigrate(&model.Order{}, &model.OrderItem{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestOrderRepository_Create(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		OrderNo:     "ORD202301010001",
		UserID:      1,
		TotalAmount: 199.99,
		Status:      model.OrderStatusPending,
	}

	err := repo.Create(order)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	if order.ID == 0 {
		t.Error("expected non-zero order ID")
	}
}

func TestOrderRepository_GetByID(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		OrderNo:     "ORD202301010001",
		UserID:      1,
		TotalAmount: 199.99,
		Status:      model.OrderStatusPending,
	}
	repo.Create(order)

	found, err := repo.GetByID(order.ID)
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}

	if found.OrderNo != order.OrderNo {
		t.Errorf("expected order_no %s, got %s", order.OrderNo, found.OrderNo)
	}
}

func TestOrderRepository_GetByOrderNo(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		OrderNo:     "ORD202301010001",
		UserID:      1,
		TotalAmount: 199.99,
		Status:      model.OrderStatusPending,
	}
	repo.Create(order)

	found, err := repo.GetByOrderNo("ORD202301010001")
	if err != nil {
		t.Fatalf("failed to get order by order_no: %v", err)
	}

	if found.TotalAmount != order.TotalAmount {
		t.Errorf("expected total_amount %f, got %f", order.TotalAmount, found.TotalAmount)
	}
}

func TestOrderRepository_Update(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		OrderNo:     "ORD202301010001",
		UserID:      1,
		TotalAmount: 199.99,
		Status:      model.OrderStatusPending,
	}
	repo.Create(order)

	order.Status = model.OrderStatusConfirmed
	err := repo.Update(order)
	if err != nil {
		t.Fatalf("failed to update order: %v", err)
	}

	found, _ := repo.GetByID(order.ID)
	if found.Status != model.OrderStatusConfirmed {
		t.Errorf("expected status 'confirmed', got '%s'", found.Status)
	}
}

func TestOrderRepository_Delete(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		OrderNo:     "ORD202301010001",
		UserID:      1,
		TotalAmount: 199.99,
		Status:      model.OrderStatusPending,
	}
	repo.Create(order)

	err := repo.Delete(order.ID)
	if err != nil {
		t.Fatalf("failed to delete order: %v", err)
	}

	_, err = repo.GetByID(order.ID)
	if err == nil {
		t.Error("expected error when getting deleted order")
	}
}

func TestOrderRepository_List(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	for i := 0; i < 5; i++ {
		repo.Create(&model.Order{
			OrderNo:     fmt.Sprintf("ORD20230101%04d", i),
			UserID:      1,
			TotalAmount: 99.99,
			Status:      model.OrderStatusPending,
		})
	}

	orders, total, err := repo.List(0, 10)
	if err != nil {
		t.Fatalf("failed to list orders: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(orders) != 5 {
		t.Errorf("expected 5 orders, got %d", len(orders))
	}
}

func TestOrderRepository_ListByUserID(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	// Create orders for different users
	repo.Create(&model.Order{
		OrderNo:     "ORD001",
		UserID:      1,
		TotalAmount: 99.99,
		Status:      model.OrderStatusPending,
	})
	repo.Create(&model.Order{
		OrderNo:     "ORD002",
		UserID:      1,
		TotalAmount: 149.99,
		Status:      model.OrderStatusPending,
	})
	repo.Create(&model.Order{
		OrderNo:     "ORD003",
		UserID:      2,
		TotalAmount: 199.99,
		Status:      model.OrderStatusPending,
	})

	orders, total, err := repo.ListByUserID(1, 0, 10)
	if err != nil {
		t.Fatalf("failed to list orders by user: %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2 for user 1, got %d", total)
	}

	if len(orders) != 2 {
		t.Errorf("expected 2 orders for user 1, got %d", len(orders))
	}
}

func TestOrderRepository_ListByUserID_Empty(t *testing.T) {
	db := setupOrderTestDB(t)
	repo := NewOrderRepository(db)

	orders, total, err := repo.ListByUserID(999, 0, 10)
	if err != nil {
		t.Fatalf("failed to list orders: %v", err)
	}

	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}

	if len(orders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(orders))
	}
}
