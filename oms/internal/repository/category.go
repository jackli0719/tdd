package repository

import (
	"oms/internal/model"

	"gorm.io/gorm"
)

// CategoryRepository handles category data access
type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id int64) (*model.Category, error)
	GetByName(name string) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id int64) error
	List(offset, limit int) ([]*model.Category, int64, error)
	HasProducts(categoryID int64) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id int64) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id int64) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *categoryRepository) List(offset, limit int) ([]*model.Category, int64, error) {
	var categories []*model.Category
	var total int64

	r.db.Model(&model.Category{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Order("id DESC").Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *categoryRepository) HasProducts(categoryID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
