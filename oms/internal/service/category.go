package service

import (
	"errors"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryHasProducts   = errors.New("category has associated products")
)

// CategoryService handles category business logic
type CategoryService interface {
	Create(req *model.CreateCategoryRequest) (*model.Category, error)
	GetByID(id int64) (*model.Category, error)
	Update(id int64, req *model.UpdateCategoryRequest) (*model.Category, error)
	Delete(id int64) error
	List(page, pageSize int) ([]*model.Category, int64, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(req *model.CreateCategoryRequest) (*model.Category, error) {
	// Check if name exists
	if _, err := s.repo.GetByName(req.Name); err == nil {
		return nil, ErrCategoryAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetByID(id int64) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(id int64, req *model.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Check if new name conflicts
	if req.Name != "" && req.Name != category.Name {
		if existing, err := s.repo.GetByName(req.Name); err == nil && existing.ID != id {
			return nil, ErrCategoryAlreadyExists
		}
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	// nil Description means clear, empty string means clear
	if req.Description != nil {
		category.Description = *req.Description
	}

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id int64) error {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	// Check if category has products
	hasProducts, err := s.repo.HasProducts(id)
	if err != nil {
		return err
	}
	if hasProducts {
		return ErrCategoryHasProducts
	}

	return s.repo.Delete(category.ID)
}

func (s *categoryService) List(page, pageSize int) ([]*model.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
