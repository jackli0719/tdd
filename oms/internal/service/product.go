package service

import (
	"errors"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductExists     = errors.New("product already exists")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// ProductService handles product business logic
type ProductService interface {
	Create(req *model.CreateProductRequest) (*model.Product, error)
	GetByID(id int64) (*model.Product, error)
	Update(id int64, req *model.UpdateProductRequest) (*model.Product, error)
	Delete(id int64) error
	List(page, pageSize int, categoryID int64) ([]*model.Product, int64, error)
	DecrementStock(id int64, quantity int) error
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new ProductService
func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{repo: repo, categoryRepo: categoryRepo}
}

func (s *productService) Create(req *model.CreateProductRequest) (*model.Product, error) {
	// Check if product with same name exists
	if _, err := s.repo.GetByName(req.Name); err == nil {
		return nil, ErrProductExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if category exists
	if _, err := s.categoryRepo.GetByID(req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	product := &model.Product{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetByID(id int64) (*model.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *productService) Update(id int64, req *model.UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Check if new name conflicts
	if req.Name != "" && req.Name != product.Name {
		if existing, err := s.repo.GetByName(req.Name); err == nil && existing.ID != id {
			return nil, ErrProductExists
		}
	}

	// Check if category exists if being updated
	if req.CategoryID != nil && *req.CategoryID != product.CategoryID {
		if _, err := s.categoryRepo.GetByID(*req.CategoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
	}

	// Update fields
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Delete(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *productService) List(page, pageSize int, categoryID int64) ([]*model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize, categoryID)
}

func (s *productService) DecrementStock(id int64, quantity int) error {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	if product.Stock < quantity {
		return ErrInsufficientStock
	}

	product.Stock -= quantity
	return s.repo.Update(product)
}
