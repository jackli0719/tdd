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

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	svc service.ProductService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// List handles GET /api/products
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseInt(c.DefaultQuery("category_id", "0"), 10, 64)

	products, total, err := h.svc.List(page, pageSize, categoryID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
	})
}

// Get handles GET /api/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	product, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, product)
}

// Create handles POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	product, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrProductExists) {
			response.Error(c, http.StatusConflict, "product already exists")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, product)
}

// Update handles PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	product, err := h.svc.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		if errors.Is(err, service.ErrProductExists) {
			response.Error(c, http.StatusConflict, "product already exists")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, product)
}

// Delete handles DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
