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

	staffs, total, err := h.svc.List(page, pageSize)
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
	var req model.Staff
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	if err := h.svc.Create(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, req)
}

// Update handles PUT /api/staff/:id
func (h *StaffHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.Staff
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	req.ID = id
	if err := h.svc.Update(&req); err != nil {
		if err == service.ErrStaffNotFound {
			response.Error(c, http.StatusNotFound, "staff not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, req)
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