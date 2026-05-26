package router

import (
	"oms/internal/handler"
	"oms/internal/middleware"
	"oms/internal/repository"
	"oms/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup configures all routes
func Setup(r *gin.Engine, db *gorm.DB) {
	// Recovery middleware must be first to catch panics
	r.Use(middleware.Recovery())

	// Logger middleware
	r.Use(middleware.Logger())

	// CORS middleware
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		if db == nil {
			c.JSON(503, gin.H{"status": "database not connected"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Skip API routes if database is not connected
	if db == nil {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(503, gin.H{"code": 503, "message": "database not connected"})
		})
		return
	}

	// API routes
	api := r.Group("/api")
	{
		// Category routes
		categoryRepo := repository.NewCategoryRepository(db)
		categorySvc := service.NewCategoryService(categoryRepo)
		categoryHandler := handler.NewCategoryHandler(categorySvc)
		api.GET("/categories", categoryHandler.List)
		api.GET("/categories/:id", categoryHandler.Get)
		api.POST("/categories", categoryHandler.Create)
		api.PUT("/categories/:id", categoryHandler.Update)
		api.DELETE("/categories/:id", categoryHandler.Delete)

		// User routes
		userRepo := repository.NewUserRepository(db)
		userSvc := service.NewUserService(userRepo)
		userHandler := handler.NewUserHandler(userSvc)
		api.GET("/users", userHandler.List)
		api.GET("/users/:id", userHandler.Get)
		api.POST("/users", userHandler.Create)
		api.PUT("/users/:id", userHandler.Update)
		api.DELETE("/users/:id", userHandler.Delete)

		// Product routes
		productRepo := repository.NewProductRepository(db)
		productSvc := service.NewProductService(productRepo, categoryRepo)
		productHandler := handler.NewProductHandler(productSvc)
		api.GET("/products", productHandler.List)
		api.GET("/products/:id", productHandler.Get)
		api.POST("/products", productHandler.Create)
		api.PUT("/products/:id", productHandler.Update)
		api.DELETE("/products/:id", productHandler.Delete)

		// Order routes
		orderRepo := repository.NewOrderRepository(db)
		orderSvc := service.NewOrderService(orderRepo, userRepo, productRepo)
		orderHandler := handler.NewOrderHandler(orderSvc)
		api.GET("/orders", orderHandler.List)
		api.GET("/orders/:id", orderHandler.Get)
		api.POST("/orders", orderHandler.Create)
		api.DELETE("/orders/:id", orderHandler.Delete)
		api.POST("/orders/:id/confirm", orderHandler.Paid)
		api.POST("/orders/:id/start", orderHandler.Ship)
		api.POST("/orders/:id/complete", orderHandler.Complete)
		api.POST("/orders/:id/cancel", orderHandler.Cancel)

		// Stats routes
		statsHandler := handler.NewStatsHandler(orderSvc)
		api.GET("/stats/orders", statsHandler.OrderStats)
		api.GET("/stats/revenue", statsHandler.RevenueStats)
	}
}
