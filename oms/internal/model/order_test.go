package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOrderJSONSerialization(t *testing.T) {
	now := time.Now()
	order := Order{
		ID:          1,
		OrderNo:     "ORD202301010001",
		UserID:      1,
		TotalAmount: 199.99,
		Status:      OrderStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal order: %v", err)
	}

	var unmarshaled Order
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal order: %v", err)
	}

	if unmarshaled.ID != order.ID {
		t.Errorf("expected ID %d, got %d", order.ID, unmarshaled.ID)
	}
	if unmarshaled.OrderNo != order.OrderNo {
		t.Errorf("expected order_no %s, got %s", order.OrderNo, unmarshaled.OrderNo)
	}
}

func TestOrderTableName(t *testing.T) {
	order := Order{}
	if order.TableName() != "orders" {
		t.Errorf("expected table name 'orders', got '%s'", order.TableName())
	}
}

func TestOrderItemTableName(t *testing.T) {
	item := OrderItem{}
	if item.TableName() != "order_items" {
		t.Errorf("expected table name 'order_items', got '%s'", item.TableName())
	}
}

func TestOrderStatus(t *testing.T) {
	if OrderStatusPending != "pending" {
		t.Errorf("expected 'pending', got '%s'", OrderStatusPending)
	}
	if OrderStatusConfirmed != "confirmed" {
		t.Errorf("expected 'confirmed', got '%s'", OrderStatusConfirmed)
	}
	if OrderStatusInService != "in_service" {
		t.Errorf("expected 'in_service', got '%s'", OrderStatusInService)
	}
	if OrderStatusCompleted != "completed" {
		t.Errorf("expected 'completed', got '%s'", OrderStatusCompleted)
	}
	if OrderStatusCancelled != "cancelled" {
		t.Errorf("expected 'cancelled', got '%s'", OrderStatusCancelled)
	}
}
