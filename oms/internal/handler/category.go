package handler

import (
	"errors"
	"net/http"
	"strconv"

	"oms/internal/model"
	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List handles GET /api/categories
func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	categories, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"categories": categories,
		"total":      total,
		"page":       page,
	})
}

// Get handles GET /api/categories/:id
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "category not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, category)
}

// Create handles POST /api/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	category, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryAlreadyExists) {
			response.Error(c, http.StatusConflict, "category already exists")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, category)
}

// Update handles PUT /api/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	category, err := h.svc.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "category not found")
			return
		}
		if errors.Is(err, service.ErrCategoryAlreadyExists) {
			response.Error(c, http.StatusConflict, "category already exists")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, category)
}

// Delete handles DELETE /api/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "category not found")
			return
		}
		if errors.Is(err, service.ErrCategoryHasProducts) {
			response.Error(c, http.StatusConflict, "category has associated products")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
