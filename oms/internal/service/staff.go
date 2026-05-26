package service

import (
	"errors"

	"oms/internal/model"
	"oms/internal/repository"
)

var (
	ErrStaffNotFound = errors.New("staff not found")
)

// StaffService handles staff business logic
type StaffService struct {
	staffRepo repository.StaffRepository
}

// NewStaffService creates a new StaffService
func NewStaffService(staffRepo repository.StaffRepository) *StaffService {
	return &StaffService{staffRepo: staffRepo}
}

// Create creates a new staff member
func (s *StaffService) Create(staff *model.Staff) error {
	return s.staffRepo.Create(staff)
}

// GetByID gets a staff member by ID
func (s *StaffService) GetByID(id int64) (*model.Staff, error) {
	staff, err := s.staffRepo.GetByID(id)
	if err != nil {
		return nil, ErrStaffNotFound
	}
	return staff, nil
}

// Update updates a staff member
func (s *StaffService) Update(staff *model.Staff) error {
	existing, err := s.staffRepo.GetByID(staff.ID)
	if err != nil {
		return ErrStaffNotFound
	}
	staff.CreatedAt = existing.CreatedAt
	return s.staffRepo.Update(staff)
}

// Delete deletes a staff member
func (s *StaffService) Delete(id int64) error {
	_, err := s.staffRepo.GetByID(id)
	if err != nil {
		return ErrStaffNotFound
	}
	return s.staffRepo.Delete(id)
}

// List returns a paginated list of staff members
func (s *StaffService) List(page, pageSize int) ([]model.Staff, int64, error) {
	return s.staffRepo.List(page, pageSize)
}

// ListAvailable returns all available staff members
func (s *StaffService) ListAvailable() ([]model.Staff, error) {
	return s.staffRepo.ListAvailable()
}

// UpdateStatus updates a staff member's status
func (s *StaffService) UpdateStatus(id int64, status model.StaffStatus) error {
	staff, err := s.staffRepo.GetByID(id)
	if err != nil {
		return ErrStaffNotFound
	}
	staff.Status = status
	return s.staffRepo.Update(staff)
}