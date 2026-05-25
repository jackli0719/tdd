package model

import (
	"time"
)

// Product represents a product in the system
type Product struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock     int       `json:"stock" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// CreateProductRequest is the request body for creating a product
type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required,min=1,max=100"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int     `json:"stock" binding:"required,gte=0"`
}

// UpdateProductRequest is the request body for updating a product
type UpdateProductRequest struct {
	Name  string  `json:"name" binding:"omitempty,min=1,max=100"`
	Price float64 `json:"price" binding:"omitempty,gt=0"`
	Stock int     `json:"stock" binding:"omitempty,gte=0"`
}
