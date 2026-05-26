package repository

import (
	"oms/internal/model"

	"gorm.io/gorm"
)

// StaffRepository handles staff data access
type StaffRepository interface {
	Create(staff *model.Staff) error
	GetByID(id int64) (*model.Staff, error)
	Update(staff *model.Staff) error
	Delete(id int64) error
	List(page, pageSize int) ([]model.Staff, int64, error)
	ListAvailable() ([]model.Staff, error)
}

type staffRepository struct {
	db *gorm.DB
}

// NewStaffRepository creates a new StaffRepository
func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(staff *model.Staff) error {
	return r.db.Create(staff).Error
}

func (r *staffRepository) GetByID(id int64) (*model.Staff, error) {
	var staff model.Staff
	if err := r.db.First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) Update(staff *model.Staff) error {
	return r.db.Save(staff).Error
}

func (r *staffRepository) Delete(id int64) error {
	return r.db.Delete(&model.Staff{}, id).Error
}

func (r *staffRepository) List(page, pageSize int) ([]model.Staff, int64, error) {
	var staffs []model.Staff
	var total int64

	r.db.Model(&model.Staff{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&staffs).Error; err != nil {
		return nil, 0, err
	}
	return staffs, total, nil
}

func (r *staffRepository) ListAvailable() ([]model.Staff, error) {
	var staffs []model.Staff
	if err := r.db.Where("status = ?", model.StaffStatusAvailable).Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}