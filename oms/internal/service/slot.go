package service

import (
	"fmt"
	"oms/internal/repository"
	"time"
)

// SlotService handles slot business logic
type SlotService struct {
	orderRepo repository.OrderRepository
}

// NewSlotService creates a new SlotService
func NewSlotService(orderRepo repository.OrderRepository) *SlotService {
	return &SlotService{orderRepo: orderRepo}
}

// Slot represents a time slot
type Slot struct {
	StartTime string `json:"start_time"` // e.g., "09:00"
	EndTime   string `json:"end_time"`   // e.g., "10:00"
	Available bool   `json:"available"`
}

// GetAvailableSlots returns available time slots for a given date
func (s *SlotService) GetAvailableSlots(date time.Time) ([]Slot, error) {
	// Start and end hours for service (9:00 - 18:00)
	startHour := 9
	endHour := 18

	// Get all appointments for this date
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	dayEnd := dayStart.Add(24 * time.Hour)

	orders, _, err := s.orderRepo.ListByDateRange(dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	// Build a map of booked hours
	bookedHours := make(map[int]bool)
	for _, order := range orders {
		if order.AppointmentTime != nil {
			hour := order.AppointmentTime.Hour()
			bookedHours[hour] = true
		}
	}

	// Generate slots (9:00-10:00, 10:00-11:00, ..., 17:00-18:00)
	var slots []Slot
	for hour := startHour; hour < endHour; hour++ {
		slot := Slot{
			StartTime: formatHour(hour),
			EndTime:   formatHour(hour + 1),
			Available: !bookedHours[hour],
		}
		slots = append(slots, slot)
	}

	return slots, nil
}

// formatHour converts hour int to string like "09:00"
func formatHour(hour int) string {
	return fmt.Sprintf("%02d:00", hour)
}