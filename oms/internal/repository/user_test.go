package repository

import (
	"fmt"

	"oms/internal/model"
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

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected non-zero user ID")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}
	repo.Create(user)

	found, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if found.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, found.Username)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}
	repo.Create(user)

	found, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("failed to get user by username: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, found.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}
	repo.Create(user)

	user.Phone = "0987654321"
	err := repo.Update(user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	found, _ := repo.GetByID(user.ID)
	if found.Phone != "0987654321" {
		t.Errorf("expected phone 0987654321, got %s", found.Phone)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}
	repo.Create(user)

	err := repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("expected error when getting deleted user")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	for i := 0; i < 5; i++ {
		repo.Create(&model.User{
			Username: fmt.Sprintf("testuser%d", i),
			Email:    fmt.Sprintf("test%d@example.com", i),
			Phone:    "1234567890",
		})
	}

	users, total, err := repo.List(0, 10)
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
