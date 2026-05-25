package service

import (
	"errors"

	"oms/internal/model"
	"oms/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidInput      = errors.New("invalid input")
)

// UserService handles user business logic
type UserService interface {
	Create(req *model.CreateUserRequest) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	Update(id int64, req *model.UpdateUserRequest) (*model.User, error)
	Delete(id int64) error
	List(page, pageSize int) ([]*model.User, int64, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req *model.CreateUserRequest) (*model.User, error) {
	// Check if username exists
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if email exists
	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id int64) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(id int64, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check if new username conflicts
	if req.Username != "" && req.Username != user.Username {
		if existing, err := s.repo.GetByUsername(req.Username); err == nil && existing.ID != id {
			return nil, ErrUserAlreadyExists
		}
	}

	// Check if new email conflicts
	if req.Email != "" && req.Email != user.Email {
		if existing, err := s.repo.GetByEmail(req.Email); err == nil && existing.ID != id {
			return nil, ErrUserAlreadyExists
		}
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *userService) List(page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
