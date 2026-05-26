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

// OrderHandler handles order HTTP requests
type OrderHandler struct {
	svc service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// List handles GET /api/orders
func (h *OrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
	})
}

// Get handles GET /api/orders/:id
func (h *OrderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	order, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// Create handles POST /api/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	order, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "product not found")
			return
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			response.Error(c, http.StatusBadRequest, "insufficient stock")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// Delete handles DELETE /api/orders/:id
func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order not found")
			return
		}
		if errors.Is(err, service.ErrInvalidOrderState) {
			response.Error(c, http.StatusBadRequest, "can only delete pending orders")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Paid handles POST /api/orders/:id/confirm (confirm order)
func (h *OrderHandler) Paid(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Paid(id); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order not found")
			return
		}
		if errors.Is(err, service.ErrInvalidOrderState) {
			response.Error(c, http.StatusBadRequest, "can only confirm pending orders")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Ship handles POST /api/orders/:id/start (start service)
func (h *OrderHandler) Ship(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Ship(id); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order not found")
			return
		}
		if errors.Is(err, service.ErrInvalidOrderState) {
			response.Error(c, http.StatusBadRequest, "can only start service for confirmed orders")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Complete handles POST /api/orders/:id/complete
func (h *OrderHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Complete(id); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order not found")
			return
		}
		if errors.Is(err, service.ErrInvalidOrderState) {
			response.Error(c, http.StatusBadRequest, "can only complete orders in service")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Cancel handles POST /api/orders/:id/cancel
func (h *OrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Cancel(id); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order not found")
			return
		}
		if errors.Is(err, service.ErrInvalidOrderState) {
			response.Error(c, http.StatusBadRequest, "invalid order state")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
