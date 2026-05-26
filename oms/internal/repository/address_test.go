package repository

import (
	"oms/internal/model"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupAddressTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.Address{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestAddressRepository_Create(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)

	address := &model.Address{
		UserID:   1,
		Name:     "张三",
		Phone:    "13800138001",
		Province: "广东省",
		City:     "深圳市",
		District: "南山区",
		Detail:   "某街道某号",
	}

	err := repo.Create(address)
	if err != nil {
		t.Fatalf("failed to create address: %v", err)
	}

	if address.ID == 0 {
		t.Error("expected non-zero address ID")
	}
}

func TestAddressRepository_GetByID(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)

	address := &model.Address{
		UserID:   1,
		Name:     "张三",
		Phone:    "13800138001",
		Province: "广东省",
		City:     "深圳市",
		District: "南山区",
		Detail:   "某街道某号",
	}
	repo.Create(address)

	found, err := repo.GetByID(address.ID)
	if err != nil {
		t.Fatalf("failed to get address: %v", err)
	}

	if found.Name != address.Name {
		t.Errorf("expected name %s, got %s", address.Name, found.Name)
	}
}

func TestAddressRepository_Update(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)

	address := &model.Address{
		UserID:   1,
		Name:     "张三",
		Phone:    "13800138001",
	}
	repo.Create(address)

	address.Name = "李四"
	err := repo.Update(address)
	if err != nil {
		t.Fatalf("failed to update address: %v", err)
	}

	found, _ := repo.GetByID(address.ID)
	if found.Name != "李四" {
		t.Errorf("expected name '李四', got '%s'", found.Name)
	}
}

func TestAddressRepository_Delete(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)

	address := &model.Address{
		UserID:   1,
		Name:     "张三",
		Phone:    "13800138001",
	}
	repo.Create(address)

	err := repo.Delete(address.ID)
	if err != nil {
		t.Fatalf("failed to delete address: %v", err)
	}

	_, err = repo.GetByID(address.ID)
	if err == nil {
		t.Error("expected error when getting deleted address")
	}
}

func TestAddressRepository_ListByUserID(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)

	// Create multiple addresses for user 1
	repo.Create(&model.Address{UserID: 1, Name: "地址1", Phone: "13800138001"})
	repo.Create(&model.Address{UserID: 1, Name: "地址2", Phone: "13800138002"})
	repo.Create(&model.Address{UserID: 2, Name: "地址3", Phone: "13800138003"})

	addresses, err := repo.ListByUserID(1)
	if err != nil {
		t.Fatalf("failed to list addresses: %v", err)
	}

	if len(addresses) != 2 {
		t.Errorf("expected 2 addresses for user 1, got %d", len(addresses))
	}
}

func TestAddressRepository_SetDefault(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)

	// Create two addresses for user 1
	addr1 := &model.Address{UserID: 1, Name: "地址1", Phone: "13800138001", IsDefault: false}
	addr2 := &model.Address{UserID: 1, Name: "地址2", Phone: "13800138002", IsDefault: true}
	repo.Create(addr1)
	repo.Create(addr2)

	// Set addr1 as default
	err := repo.SetDefault(addr1.ID, 1)
	if err != nil {
		t.Fatalf("failed to set default: %v", err)
	}

	// Verify addr1 is now default
	found1, _ := repo.GetByID(addr1.ID)
	if !found1.IsDefault {
		t.Error("expected addr1 to be default")
	}

	// Verify addr2 is no longer default
	found2, _ := repo.GetByID(addr2.ID)
	if found2.IsDefault {
		t.Error("expected addr2 to not be default")
	}
}