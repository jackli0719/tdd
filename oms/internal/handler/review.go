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

// ReviewHandler handles review HTTP requests.
type ReviewHandler struct {
	svc *service.ReviewService
}

// NewReviewHandler creates a new ReviewHandler.
func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// Create handles POST /api/reviews.
func (h *ReviewHandler) Create(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "未提供认证用户")
		return
	}
	userID, ok := userIDValue.(int64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "认证用户无效")
		return
	}

	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	review, err := h.svc.Create(userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"review": review})
}

// Get handles GET /api/reviews/:id.
func (h *ReviewHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	review, err := h.svc.GetByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"review": review})
}

// List handles GET /api/reviews.
func (h *ReviewHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if orderIDParam := c.Query("order_id"); orderIDParam != "" {
		orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid order_id")
			return
		}
		review, err := h.svc.GetByOrderID(orderID)
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.Success(c, gin.H{"review": review})
		return
	}

	if staffIDParam := c.Query("staff_id"); staffIDParam != "" {
		staffID, err := strconv.ParseInt(staffIDParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid staff_id")
			return
		}
		reviews, total, err := h.svc.ListByStaffID(staffID, page, pageSize)
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.Success(c, gin.H{"reviews": reviews, "total": total, "page": page})
		return
	}

	reviews, total, err := h.svc.List(page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"reviews": reviews, "total": total, "page": page})
}

// StaffSummary handles GET /api/reviews/staff-summary.
func (h *ReviewHandler) StaffSummary(c *gin.Context) {
	if staffIDParam := c.Query("staff_id"); staffIDParam != "" {
		staffID, err := strconv.ParseInt(staffIDParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid staff_id")
			return
		}
		summary, err := h.svc.GetStaffSummary(staffID)
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.Success(c, gin.H{"summary": summary})
		return
	}

	summaries, err := h.svc.ListStaffSummaries()
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, gin.H{"summaries": summaries})
}

func (h *ReviewHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrReviewNotFound):
		response.Error(c, http.StatusNotFound, "review not found")
	case errors.Is(err, service.ErrOrderNotFound):
		response.Error(c, http.StatusNotFound, "order not found")
	case errors.Is(err, service.ErrInvalidRating):
		response.Error(c, http.StatusBadRequest, "rating must be between 1 and 5")
	case errors.Is(err, service.ErrReviewAlreadyExists):
		response.Error(c, http.StatusBadRequest, "review already exists for order")
	case errors.Is(err, service.ErrReviewOrderInvalid):
		response.Error(c, http.StatusBadRequest, "only completed orders with assigned staff can be reviewed")
	default:
		response.Error(c, http.StatusInternalServerError, err.Error())
	}
}
