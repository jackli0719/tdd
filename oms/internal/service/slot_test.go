package service

import (
	"testing"
	"time"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSlotTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.Order{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestSlotService_GetAvailableSlots(t *testing.T) {
	db := setupSlotTestDB(t)
	repo := repository.NewOrderRepository(db)
	svc := NewSlotService(repo)

	// Test getting slots for a date with no orders
	tomorrow := time.Now().AddDate(0, 0, 1)
	slots, err := svc.GetAvailableSlots(tomorrow)
	if err != nil {
		t.Fatalf("GetAvailableSlots() error = %v", err)
	}

	// Should have 9 slots (9:00-18:00)
	if len(slots) != 9 {
		t.Errorf("expected 9 slots, got %d", len(slots))
	}

	// All should be available
	for _, slot := range slots {
		if !slot.Available {
			t.Errorf("expected slot %s to be available", slot.StartTime)
		}
	}
}

func TestSlotService_GetAvailableSlots_WithBookedSlot(t *testing.T) {
	db := setupSlotTestDB(t)
	repo := repository.NewOrderRepository(db)
	svc := NewSlotService(repo)

	// Create an order with appointment time tomorrow at 10:00
	tomorrow := time.Now().AddDate(0, 0, 1)
	appointmentTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, time.Local)

	order := &model.Order{
		OrderNo:         "TEST001",
		UserID:          1,
		TotalAmount:     100,
		Status:          model.OrderStatusPending,
		AppointmentTime: &appointmentTime,
	}
	repo.Create(order)

	slots, err := svc.GetAvailableSlots(tomorrow)
	if err != nil {
		t.Fatalf("GetAvailableSlots() error = %v", err)
	}

	// Should have 9 slots
	if len(slots) != 9 {
		t.Errorf("expected 9 slots, got %d", len(slots))
	}

	// Find the 10:00 slot - should be unavailable
	foundUnavailable := false
	foundAvailable := false
	for _, slot := range slots {
		if slot.StartTime == "10:00" {
			if slot.Available {
				t.Errorf("10:00 slot should be unavailable but was available")
			}
			foundUnavailable = true
		} else {
			if slot.Available {
				foundAvailable = true
			}
		}
	}

	if !foundUnavailable {
		t.Errorf("10:00 slot not found in slots")
	}
	if !foundAvailable {
		t.Errorf("some slots should be available")
	}
}

func TestSlotService_GetAvailableSlots_MultipleOrdersSameHour(t *testing.T) {
	db := setupSlotTestDB(t)
	repo := repository.NewOrderRepository(db)
	svc := NewSlotService(repo)

	// Create two orders at the same hour
	tomorrow := time.Now().AddDate(0, 0, 1)
	appointmentTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 14, 0, 0, 0, time.Local)

	order1 := &model.Order{
		OrderNo:         "TEST001",
		UserID:          1,
		TotalAmount:     100,
		Status:          model.OrderStatusPending,
		AppointmentTime: &appointmentTime,
	}
	repo.Create(order1)

	order2 := &model.Order{
		OrderNo:         "TEST002",
		UserID:          2,
		TotalAmount:     200,
		Status:          model.OrderStatusPending,
		AppointmentTime: &appointmentTime,
	}
	repo.Create(order2)

	slots, err := svc.GetAvailableSlots(tomorrow)
	if err != nil {
		t.Fatalf("GetAvailableSlots() error = %v", err)
	}

	// Find the 14:00 slot - should still just be one slot marked unavailable
	for _, slot := range slots {
		if slot.StartTime == "14:00" {
			if slot.Available {
				t.Errorf("14:00 slot should be unavailable but was available")
			}
		}
	}
}

func TestSlotService_GetAvailableSlots_SlotTimes(t *testing.T) {
	db := setupSlotTestDB(t)
	repo := repository.NewOrderRepository(db)
	svc := NewSlotService(repo)

	tomorrow := time.Now().AddDate(0, 0, 1)
	slots, err := svc.GetAvailableSlots(tomorrow)
	if err != nil {
		t.Fatalf("GetAvailableSlots() error = %v", err)
	}

	// Verify slot times: 09:00-10:00, 10:00-11:00, ..., 17:00-18:00
	expectedTimes := []string{
		"09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00",
	}

	if len(slots) != len(expectedTimes) {
		t.Errorf("expected %d slots, got %d", len(expectedTimes), len(slots))
	}

	for i, slot := range slots {
		if slot.StartTime != expectedTimes[i] {
			t.Errorf("slot[%d] StartTime = %s, want %s", i, slot.StartTime, expectedTimes[i])
		}
		if slot.EndTime != formatHour(i + 10) {
			t.Errorf("slot[%d] EndTime = %s, want %s", i, slot.EndTime, formatHour(i+10))
		}
	}
}