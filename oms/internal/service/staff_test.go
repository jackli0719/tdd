package service

import (
	"oms/internal/model"
	"oms/internal/repository"
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

func TestStaffService_Create(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	tests := []struct {
		name    string
		staff   *model.Staff
		wantErr bool
		errType error
	}{
		{
			name:    "valid staff",
			staff:   &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable},
			wantErr: false,
		},
		{
			name:    "empty name",
			staff:   &model.Staff{Name: "", Phone: "13800138001"},
			wantErr: true,
			errType: ErrEmptyName,
		},
		{
			name:    "invalid status",
			staff:   &model.Staff{Name: "李四", Phone: "13800138002", Status: "invalid"},
			wantErr: true,
			errType: ErrInvalidStatus,
		},
		{
			name:    "default status to available",
			staff:   &model.Staff{Name: "王五", Phone: "13800138003"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Create(tt.staff)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("Create() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestStaffService_Update(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	// Create a staff first
	staff := &model.Staff{Name: "初始名", Phone: "13800138000", Status: model.StaffStatusAvailable}
	if err := svc.Create(staff); err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		nameArg string
		phone   string
		status  string
		wantErr bool
		errType error
	}{
		{
			name:    "valid update",
			id:      staff.ID,
			nameArg: "新名字",
			phone:   "13900139000",
			status:  string(model.StaffStatusBusy),
			wantErr: false,
		},
		{
			name:    "empty name",
			id:      staff.ID,
			nameArg: "",
			phone:   "13900139000",
			wantErr: true,
			errType: ErrEmptyName,
		},
		{
			name:    "not found",
			id:      99999,
			nameArg: "不存在",
			phone:   "13900139000",
			wantErr: true,
			errType: ErrStaffNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Update(tt.id, tt.nameArg, tt.phone, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("Update() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestStaffService_GetByID(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	// Create a staff
	staff := &model.Staff{Name: "测试", Phone: "13800138000"}
	if err := svc.Create(staff); err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{"found", staff.ID, false},
		{"not found", 99999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffService_Delete(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	// Create a staff
	staff := &model.Staff{Name: "删除测试", Phone: "13800138000"}
	if err := svc.Create(staff); err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{"success", staff.ID, false},
		{"not found", 99999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffService_List(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	// Create multiple staff
	for i := 0; i < 5; i++ {
		staff := &model.Staff{Name: "员工" + string(rune('A'+i)), Phone: "1380013800" + string(rune('0'+i))}
		if err := svc.Create(staff); err != nil {
			t.Fatalf("failed to create staff: %v", err)
		}
	}

	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantCount  int
		wantTotal  int64
	}{
		{"first page", 1, 2, 2, 5},
		{"second page", 2, 2, 2, 5},
		{"page beyond data", 10, 10, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staffs, total, err := svc.List(tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if len(staffs) != tt.wantCount {
				t.Errorf("List() got %v staff, want %v", len(staffs), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("List() total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestStaffService_ListAvailable(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	// Create staff with different statuses
	svc.Create(&model.Staff{Name: "空闲1", Phone: "13800138001", Status: model.StaffStatusAvailable})
	svc.Create(&model.Staff{Name: "忙碌1", Phone: "13800138002", Status: model.StaffStatusBusy})
	svc.Create(&model.Staff{Name: "空闲2", Phone: "13800138003", Status: model.StaffStatusAvailable})
	svc.Create(&model.Staff{Name: "休息1", Phone: "13800138004", Status: model.StaffStatusOff})

	available, err := svc.ListAvailable()
	if err != nil {
		t.Fatalf("ListAvailable() error = %v", err)
	}
	if len(available) != 2 {
		t.Errorf("ListAvailable() got %v available staff, want 2", len(available))
	}
	for _, s := range available {
		if s.Status != model.StaffStatusAvailable {
			t.Errorf("ListAvailable() returned staff with status %v, want available", s.Status)
		}
	}
}

func TestStaffService_UpdateStatus(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := NewStaffService(repo)

	// Create a staff
	staff := &model.Staff{Name: "状态测试", Phone: "13800138000", Status: model.StaffStatusAvailable}
	if err := svc.Create(staff); err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		status  model.StaffStatus
		wantErr bool
		errType error
	}{
		{"set to busy", staff.ID, model.StaffStatusBusy, false, nil},
		{"set to off", staff.ID, model.StaffStatusOff, false, nil},
		{"set to available", staff.ID, model.StaffStatusAvailable, false, nil},
		{"invalid status", staff.ID, "invalid", true, ErrInvalidStatus},
		{"not found", 99999, model.StaffStatusAvailable, true, ErrStaffNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UpdateStatus(tt.id, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("UpdateStatus() error = %v, want %v", err, tt.errType)
			}
		})
	}
}