package repository

import (
	"oms/internal/model"

	"gorm.io/gorm"
)

// AddressRepository handles address data access
type AddressRepository interface {
	Create(address *model.Address) error
	GetByID(id int64) (*model.Address, error)
	Update(address *model.Address) error
	Delete(id int64) error
	ListByUserID(userID int64) ([]model.Address, error)
	SetDefault(id, userID int64) error
}

type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new AddressRepository
func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) GetByID(id int64) (*model.Address, error) {
	var address model.Address
	if err := r.db.First(&address, id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id int64) error {
	return r.db.Delete(&model.Address{}, id).Error
}

func (r *addressRepository) ListByUserID(userID int64) ([]model.Address, error) {
	var addresses []model.Address
	if err := r.db.Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) SetDefault(id, userID int64) error {
	tx := r.db.Begin()

	// Unset all defaults for this user
	if err := tx.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the specified address as default
	if err := tx.Model(&model.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}