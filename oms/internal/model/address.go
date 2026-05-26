package model

import (
	"time"
)

// Address represents a user's service address
type Address struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null"`
	Phone     string    `json:"phone" gorm:"type:varchar(20);not null"`
	Province  string    `json:"province" gorm:"type:varchar(50)"`
	City      string    `json:"city" gorm:"type:varchar(50)"`
	District  string    `json:"district" gorm:"type:varchar(50)"`
	Detail    string    `json:"detail" gorm:"type:varchar(255)"`
	IsDefault bool      `json:"is_default" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Address
func (Address) TableName() string {
	return "addresses"
}

// CreateAddressRequest is the request body for creating an address
type CreateAddressRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Name     string `json:"name" binding:"required,min=1,max=50"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Province string `json:"province" binding:"omitempty,max=50"`
	City     string `json:"city" binding:"omitempty,max=50"`
	District string `json:"district" binding:"omitempty,max=50"`
	Detail   string `json:"detail" binding:"omitempty,max=255"`
}

// UpdateAddressRequest is the request body for updating an address
type UpdateAddressRequest struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=50"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Province string `json:"province" binding:"omitempty,max=50"`
	City     string `json:"city" binding:"omitempty,max=50"`
	District string `json:"district" binding:"omitempty,max=50"`
	Detail   string `json:"detail" binding:"omitempty,max=255"`
}