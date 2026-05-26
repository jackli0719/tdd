package service

import (
	"errors"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrAddressNotFound = errors.New("address not found")
)

// AddressService handles address business logic
type AddressService struct {
	addressRepo repository.AddressRepository
}

// NewAddressService creates a new AddressService
func NewAddressService(addressRepo repository.AddressRepository) *AddressService {
	return &AddressService{addressRepo: addressRepo}
}

// Create creates a new address
func (s *AddressService) Create(req *model.CreateAddressRequest) (*model.Address, error) {
	// If this is the first address for the user, set it as default
	addresses, err := s.addressRepo.ListByUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	address := &model.Address{
		UserID:    req.UserID,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		IsDefault: len(addresses) == 0,
	}

	if err := s.addressRepo.Create(address); err != nil {
		return nil, err
	}
	return address, nil
}

// GetByID gets an address by ID
func (s *AddressService) GetByID(id int64) (*model.Address, error) {
	address, err := s.addressRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

// Update updates an address
func (s *AddressService) Update(id int64, req *model.UpdateAddressRequest) (*model.Address, error) {
	address, err := s.addressRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		address.Name = req.Name
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.Province != "" {
		address.Province = req.Province
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.District != "" {
		address.District = req.District
	}
	if req.Detail != "" {
		address.Detail = req.Detail
	}

	if err := s.addressRepo.Update(address); err != nil {
		return nil, err
	}
	return address, nil
}

// Delete deletes an address
func (s *AddressService) Delete(id int64) error {
	_, err := s.addressRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAddressNotFound
		}
		return err
	}
	return s.addressRepo.Delete(id)
}

// ListByUserID lists all addresses for a user
func (s *AddressService) ListByUserID(userID int64) ([]model.Address, error) {
	return s.addressRepo.ListByUserID(userID)
}

// SetDefault sets an address as the default
func (s *AddressService) SetDefault(id, userID int64) error {
	_, err := s.addressRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAddressNotFound
		}
		return err
	}
	return s.addressRepo.SetDefault(id, userID)
}