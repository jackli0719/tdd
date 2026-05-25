package service

import (
	"fmt"

	"oms/internal/model"
	"oms/internal/repository"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestUserService_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}

	user, err := svc.Create(req)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("expected username %s, got %s", req.Username, user.Username)
	}
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}

	svc.Create(req)

	// Try to create another user with same username
	_, err := svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test2@example.com",
		Phone:    "0987654321",
	})

	if err != ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestUserService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}
	created, _ := svc.Create(req)

	found, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if found.Username != created.Username {
		t.Errorf("expected username %s, got %s", created.Username, found.Username)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	_, err := svc.GetByID(999)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	created, _ := svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	})

	updated, err := svc.Update(created.ID, &model.UpdateUserRequest{
		Phone: "0987654321",
	})
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	if updated.Phone != "0987654321" {
		t.Errorf("expected phone 0987654321, got %s", updated.Phone)
	}
}

func TestUserService_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	created, _ := svc.Create(&model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	})

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	_, err = svc.GetByID(created.ID)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_List(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo)

	for i := 0; i < 5; i++ {
		svc.Create(&model.CreateUserRequest{
			Username: fmt.Sprintf("testuser%d", i),
			Email:    fmt.Sprintf("test%d@example.com", i),
			Phone:    "1234567890",
		})
	}

	users, total, err := svc.List(1, 10)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(users) != 5 {
		t.Errorf("expected 5 users, got %d", len(users))
	}
}
