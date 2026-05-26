package handler

import (
	"net/http"
	"time"

	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
)

// SlotHandler handles slot HTTP requests
type SlotHandler struct {
	svc *service.SlotService
}

// NewSlotHandler creates a new SlotHandler
func NewSlotHandler(svc *service.SlotService) *SlotHandler {
	return &SlotHandler{svc: svc}
}

// List handles GET /api/slots?date=2026-05-26
func (h *SlotHandler) List(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		response.Error(c, http.StatusBadRequest, "date parameter is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	slots, err := h.svc.GetAvailableSlots(date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"slots": slots,
		"date":  dateStr,
	})
}