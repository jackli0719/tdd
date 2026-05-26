package repository

import (
	"oms/internal/model"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupStaffTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&model.Staff{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestStaffRepository_Create(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	staff := &model.Staff{
		Name:   "张三",
		Phone:  "13800138001",
		Status: model.StaffStatusAvailable,
	}

	err := repo.Create(staff)
	if err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	if staff.ID == 0 {
		t.Error("expected non-zero staff ID")
	}
}

func TestStaffRepository_GetByID(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	staff := &model.Staff{
		Name:   "张三",
		Phone:  "13800138001",
		Status: model.StaffStatusAvailable,
	}
	repo.Create(staff)

	found, err := repo.GetByID(staff.ID)
	if err != nil {
		t.Fatalf("failed to get staff: %v", err)
	}

	if found.Name != staff.Name {
		t.Errorf("expected name %s, got %s", staff.Name, found.Name)
	}
}

func TestStaffRepository_GetByID_NotFound(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("expected error when getting non-existent staff")
	}
}

func TestStaffRepository_Update(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	staff := &model.Staff{
		Name:   "张三",
		Phone:  "13800138001",
		Status: model.StaffStatusAvailable,
	}
	repo.Create(staff)

	staff.Phone = "13900139001"
	staff.Status = model.StaffStatusBusy
	err := repo.Update(staff)
	if err != nil {
		t.Fatalf("failed to update staff: %v", err)
	}

	found, _ := repo.GetByID(staff.ID)
	if found.Phone != "13900139001" {
		t.Errorf("expected phone 13900139001, got %s", found.Phone)
	}
	if found.Status != model.StaffStatusBusy {
		t.Errorf("expected status busy, got %s", found.Status)
	}
}

func TestStaffRepository_Delete(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	staff := &model.Staff{
		Name:   "张三",
		Phone:  "13800138001",
		Status: model.StaffStatusAvailable,
	}
	repo.Create(staff)

	err := repo.Delete(staff.ID)
	if err != nil {
		t.Fatalf("failed to delete staff: %v", err)
	}

	_, err = repo.GetByID(staff.ID)
	if err == nil {
		t.Error("expected error when getting deleted staff")
	}
}

func TestStaffRepository_List(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	for i := 0; i < 5; i++ {
		repo.Create(&model.Staff{
			Name:   "员工",
			Phone:  "1380013800",
			Status: model.StaffStatusAvailable,
		})
	}

	staffs, total, err := repo.List(1, 3)
	if err != nil {
		t.Fatalf("failed to list staff: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(staffs) != 3 {
		t.Errorf("expected 3 staff, got %d", len(staffs))
	}
}

func TestStaffRepository_List_Pagination(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	for i := 0; i < 5; i++ {
		repo.Create(&model.Staff{
			Name:   "员工",
			Phone:  "1380013800",
			Status: model.StaffStatusAvailable,
		})
	}

	// Second page
	staffs, total, err := repo.List(2, 3)
	if err != nil {
		t.Fatalf("failed to list staff page 2: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(staffs) != 2 {
		t.Errorf("expected 2 staff on page 2, got %d", len(staffs))
	}
}

func TestStaffRepository_ListAvailable(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	repo.Create(&model.Staff{Name: "空闲1", Phone: "13800138001", Status: model.StaffStatusAvailable})
	repo.Create(&model.Staff{Name: "忙碌1", Phone: "13800138002", Status: model.StaffStatusBusy})
	repo.Create(&model.Staff{Name: "空闲2", Phone: "13800138003", Status: model.StaffStatusAvailable})

	available, err := repo.ListAvailable()
	if err != nil {
		t.Fatalf("failed to list available staff: %v", err)
	}

	if len(available) != 2 {
		t.Errorf("expected 2 available staff, got %d", len(available))
	}

	for _, s := range available {
		if s.Status != model.StaffStatusAvailable {
			t.Errorf("expected available status, got %s", s.Status)
		}
	}
}

func TestStaffRepository_ListByStatus(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	repo.Create(&model.Staff{Name: "空闲1", Phone: "13800138001", Status: model.StaffStatusAvailable})
	repo.Create(&model.Staff{Name: "忙碌1", Phone: "13800138002", Status: model.StaffStatusBusy})
	repo.Create(&model.Staff{Name: "空闲2", Phone: "13800138003", Status: model.StaffStatusAvailable})

	staffs, total, err := repo.ListByStatus(model.StaffStatusAvailable, 1, 10)
	if err != nil {
		t.Fatalf("failed to list by status: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 available staff, got %d", total)
	}

	if len(staffs) != 2 {
		t.Errorf("expected 2 staff, got %d", len(staffs))
	}

	for _, s := range staffs {
		if s.Status != model.StaffStatusAvailable {
			t.Errorf("expected available status, got %s", s.Status)
		}
	}
}

func TestStaffRepository_ListByStatus_Pagination(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := NewStaffRepository(db)

	// Create 3 available staff
	for i := 0; i < 3; i++ {
		repo.Create(&model.Staff{
			Name:   "员工",
			Phone:  "1380013800",
			Status: model.StaffStatusAvailable,
		})
	}
	repo.Create(&model.Staff{Name: "忙碌", Phone: "13800138002", Status: model.StaffStatusBusy})

	// List only first 2 available
	staffs, total, err := repo.ListByStatus(model.StaffStatusAvailable, 1, 2)
	if err != nil {
		t.Fatalf("failed to list by status: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3 available, got %d", total)
	}

	if len(staffs) != 2 {
		t.Errorf("expected 2 staff, got %d", len(staffs))
	}
}