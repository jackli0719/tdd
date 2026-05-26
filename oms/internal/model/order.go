package model

import (
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusInService  OrderStatus = "in_service"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID          int64       `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo     string      `json:"order_no" gorm:"type:varchar(32);not null;uniqueIndex:uk_order_no"`
	UserID      int64       `json:"user_id" gorm:"not null;index:idx_user_id"`
	TotalAmount float64     `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(20);not null;default:pending;index:idx_status"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int64     `json:"order_id" gorm:"not null;index:idx_order_id"`
	ProductID int64     `json:"product_id" gorm:"not null;index:idx_product_id"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Subtotal  float64   `json:"subtotal" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specifies the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// CreateOrderRequest is the request body for creating an order
type CreateOrderRequest struct {
	UserID int64                    `json:"user_id" binding:"required"`
	Items  []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

// CreateOrderItemRequest is the request body for an order item
type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

// OrderResponse is the response for an order with items
type OrderResponse struct {
	Order
	Items []OrderItem `json:"items"`
}
