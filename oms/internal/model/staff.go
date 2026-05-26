package model

import (
	"time"
)

// StaffStatus represents the status of a staff member
type StaffStatus string

const (
	StaffStatusAvailable StaffStatus = "available" // 空闲
	StaffStatusBusy      StaffStatus = "busy"      // 忙碌
	StaffStatusOff       StaffStatus = "off"       // 休息
)

// Staff represents a service staff member
type Staff struct {
	ID        int64       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string      `json:"name" gorm:"type:varchar(50);not null"`
	Phone     string      `json:"phone" gorm:"type:varchar(20)"`
	Avatar    string      `json:"avatar" gorm:"type:varchar(255)"`
	Status    StaffStatus `json:"status" gorm:"type:varchar(20);not null;default:available;index"`
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Staff
func (Staff) TableName() string {
	return "staff"
}

// CreateStaffRequest is the request body for creating a staff member
type CreateStaffRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=50"`
	Phone  string `json:"phone" binding:"omitempty,max=20"`
	Avatar string `json:"avatar" binding:"omitempty,max=255"`
	Status string `json:"status" binding:"omitempty,oneof=available busy off"`
}

// UpdateStaffRequest is the request body for updating a staff member
type UpdateStaffRequest struct {
	ID     int64  `json:"id"`
	Name   string `json:"name" binding:"omitempty,min=1,max=50"`
	Phone  string `json:"phone" binding:"omitempty,max=20"`
	Avatar string `json:"avatar" binding:"omitempty,max=255"`
	Status string `json:"status" binding:"omitempty,oneof=available busy off"`
}
