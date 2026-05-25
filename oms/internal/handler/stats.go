package handler

import (
	"net/http"

	"oms/internal/model"
	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatsHandler handles stats HTTP requests
type StatsHandler struct {
	orderSvc service.OrderService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(orderSvc service.OrderService) *StatsHandler {
	return &StatsHandler{orderSvc: orderSvc}
}

// OrderStats represents order statistics
type OrderStats struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	Paid      int64 `json:"paid"`
	Shipped   int64 `json:"shipped"`
	Completed int64 `json:"completed"`
	Cancelled int64 `json:"cancelled"`
}

// RevenueStats represents revenue statistics
type RevenueStats struct {
	TotalRevenue     float64 `json:"total_revenue"`
	PendingRevenue  float64 `json:"pending_revenue"`
	PaidRevenue     float64 `json:"paid_revenue"`
	ShippedRevenue  float64 `json:"shipped_revenue"`
	CompletedRevenue float64 `json:"completed_revenue"`
}

// OrderStats handles GET /api/stats/orders
func (h *StatsHandler) OrderStats(c *gin.Context) {
	orders, total, err := h.orderSvc.List(1, 1000)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	stats := OrderStats{Total: total}
	for _, order := range orders {
		switch order.Status {
		case model.OrderStatusPending:
			stats.Pending++
		case model.OrderStatusPaid:
			stats.Paid++
		case model.OrderStatusShipped:
			stats.Shipped++
		case model.OrderStatusCompleted:
			stats.Completed++
		case model.OrderStatusCancelled:
			stats.Cancelled++
		}
	}

	response.Success(c, stats)
}

// RevenueStats handles GET /api/stats/revenue
func (h *StatsHandler) RevenueStats(c *gin.Context) {
	orders, total, err := h.orderSvc.List(1, 1000)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	stats := RevenueStats{}
	for i := int64(0); i < total; i++ {
		order := orders[i]
		stats.TotalRevenue += order.TotalAmount

		switch order.Status {
		case model.OrderStatusPending:
			stats.PendingRevenue += order.TotalAmount
		case model.OrderStatusPaid:
			stats.PaidRevenue += order.TotalAmount
		case model.OrderStatusShipped:
			stats.ShippedRevenue += order.TotalAmount
		case model.OrderStatusCompleted:
			stats.CompletedRevenue += order.TotalAmount
		}
	}

	response.Success(c, stats)
}
