package handler

import (
	"net/http"
	"strconv"

	"oms/internal/model"
	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
)

// StaffHandler handles staff HTTP requests
type StaffHandler struct {
	svc *service.StaffService
}

// NewStaffHandler creates a new StaffHandler
func NewStaffHandler(svc *service.StaffService) *StaffHandler {
	return &StaffHandler{svc: svc}
}

// List handles GET /api/staff
func (h *StaffHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var staffs []model.Staff
	var total int64
	var err error

	if status != "" {
		// Validate status enum
		validStatuses := []model.StaffStatus{model.StaffStatusAvailable, model.StaffStatusBusy, model.StaffStatusOff}
		valid := false
		for _, s := range validStatuses {
			if model.StaffStatus(status) == s {
				valid = true
				break
			}
		}
		if !valid {
			response.Error(c, http.StatusBadRequest, "invalid status: must be available, busy, or off")
			return
		}
		staffs, total, err = h.svc.ListByStatus(model.StaffStatus(status), page, pageSize)
	} else {
		var staffList []model.Staff
		staffList, total, err = h.svc.List(page, pageSize)
		staffs = staffList
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"staffs": staffs,
		"total":  total,
		"page":   page,
	})
}

// Get handles GET /api/staff/:id
func (h *StaffHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	staff, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "staff not found")
		return
	}

	response.Success(c, staff)
}

// Create handles POST /api/staff
func (h *StaffHandler) Create(c *gin.Context) {
	var req model.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	staff := &model.Staff{
		Name:   req.Name,
		Phone:  req.Phone,
		Avatar: req.Avatar,
		Status: model.StaffStatus(req.Status),
	}
	if staff.Status == "" {
		staff.Status = model.StaffStatusAvailable
	}

	if err := h.svc.Create(staff); err != nil {
		if err == service.ErrEmptyName {
			response.Error(c, http.StatusBadRequest, "name cannot be empty")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "invalid status: must be available, busy, or off")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, staff)
}

// Update handles PUT /api/staff/:id
func (h *StaffHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	req.ID = id
	if err := h.svc.Update(req.ID, req.Name, req.Phone, req.Avatar, req.Status); err != nil {
		if err == service.ErrStaffNotFound {
			response.Error(c, http.StatusNotFound, "staff not found")
			return
		}
		if err == service.ErrEmptyName {
			response.Error(c, http.StatusBadRequest, "name cannot be empty")
			return
		}
		if err == service.ErrInvalidStatus {
			response.Error(c, http.StatusBadRequest, "invalid status: must be available, busy, or off")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete handles DELETE /api/staff/:id
func (h *StaffHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if err == service.ErrStaffNotFound {
			response.Error(c, http.StatusNotFound, "staff not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus handles PUT /api/staff/:id/status
func (h *StaffHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Status model.StaffStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// Validate status enum
	switch req.Status {
	case model.StaffStatusAvailable, model.StaffStatusBusy, model.StaffStatusOff:
		// Valid
	default:
		response.Error(c, http.StatusBadRequest, "invalid status: must be available, busy, or off")
		return
	}

	if err := h.svc.UpdateStatus(id, req.Status); err != nil {
		if err == service.ErrStaffNotFound {
			response.Error(c, http.StatusNotFound, "staff not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
