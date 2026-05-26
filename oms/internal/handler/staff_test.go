package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"oms/internal/model"
	"oms/internal/repository"
	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
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

func setupStaffRouter(svc *service.StaffService) *gin.Engine {
	r := gin.New()
	h := NewStaffHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/staff", h.List)
		api.GET("/staff/:id", h.Get)
		api.POST("/staff", h.Create)
		api.PUT("/staff/:id", h.Update)
		api.DELETE("/staff/:id", h.Delete)
		api.PUT("/staff/:id/status", h.UpdateStatus)
	}

	return r
}

func TestStaffHandler_Create(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	body := `{"name":"张三","phone":"13800138001"}`
	req, _ := http.NewRequest("POST", "/api/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
}

func TestStaffHandler_Create_InvalidJSON(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	body := `{invalid json}`
	req, _ := http.NewRequest("POST", "/api/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestStaffHandler_Create_EmptyName(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	body := `{"name":"","phone":"13800138001"}`
	req, _ := http.NewRequest("POST", "/api/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestStaffHandler_List(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	// Create a staff first
	svc.Create(&model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable})

	req, _ := http.NewRequest("GET", "/api/staff?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	staffs := data["staffs"].([]interface{})
	if len(staffs) != 1 {
		t.Errorf("expected 1 staff, got %d", len(staffs))
	}
}

func TestStaffHandler_List_FilterByStatus(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	// Create staff with different statuses
	svc.Create(&model.Staff{Name: "空闲", Phone: "13800138001", Status: model.StaffStatusAvailable})
	svc.Create(&model.Staff{Name: "忙碌", Phone: "13800138002", Status: model.StaffStatusBusy})

	// Filter by available
	req, _ := http.NewRequest("GET", "/api/staff?status=available", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	staffs := data["staffs"].([]interface{})
	if len(staffs) != 1 {
		t.Errorf("expected 1 available staff, got %d", len(staffs))
	}
}

func TestStaffHandler_List_FilterByInvalidStatus(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	req, _ := http.NewRequest("GET", "/api/staff?status=invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestStaffHandler_Get(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	staff := &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable}
	svc.Create(staff)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/staff/%d", staff.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStaffHandler_Get_NotFound(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	req, _ := http.NewRequest("GET", "/api/staff/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestStaffHandler_Update(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	staff := &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable}
	svc.Create(staff)

	body := `{"name":"李四","phone":"13900139001"}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/staff/%d", staff.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStaffHandler_Update_NotFound(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	body := `{"name":"李四"}`
	req, _ := http.NewRequest("PUT", "/api/staff/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestStaffHandler_Update_EmptyName(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	staff := &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable}
	svc.Create(staff)

	body := `{"name":""}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/staff/%d", staff.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestStaffHandler_Delete(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	staff := &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable}
	svc.Create(staff)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/staff/%d", staff.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStaffHandler_Delete_NotFound(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	req, _ := http.NewRequest("DELETE", "/api/staff/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestStaffHandler_UpdateStatus(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	staff := &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable}
	svc.Create(staff)

	body := `{"status":"busy"}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/staff/%d/status", staff.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStaffHandler_UpdateStatus_InvalidStatus(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	staff := &model.Staff{Name: "张三", Phone: "13800138001", Status: model.StaffStatusAvailable}
	svc.Create(staff)

	body := `{"status":"invalid"}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/staff/%d/status", staff.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestStaffHandler_UpdateStatus_NotFound(t *testing.T) {
	db := setupStaffTestDB(t)
	repo := repository.NewStaffRepository(db)
	svc := service.NewStaffService(repo)
	r := setupStaffRouter(svc)

	body := `{"status":"busy"}`
	req, _ := http.NewRequest("PUT", "/api/staff/999/status", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}