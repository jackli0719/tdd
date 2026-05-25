package repository

import (
	"oms/internal/model"

	"gorm.io/gorm"
)

// ProductRepository handles product data access
type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id int64) (*model.Product, error)
	GetByName(name string) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id int64) error
	List(offset, limit int) ([]*model.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id int64) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByName(name string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("name = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id int64) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) List(offset, limit int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	r.db.Model(&model.Product{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Order("id DESC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
