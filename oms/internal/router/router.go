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
func Setup(r *gin.Engine, db *gorm.DB, jwtSecret string) {
	// Recovery middleware must be first to catch panics
	r.Use(middleware.Recovery())

	// Logger middleware
	r.Use(middleware.Logger())

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

	// Repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	staffRepo := repository.NewStaffRepository(db)

	// Services
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, jwtSecret)
	productSvc := service.NewProductService(productRepo)
	orderSvc := service.NewOrderService(orderRepo, userRepo, productRepo, staffRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	statsHandler := handler.NewStatsHandler(orderSvc)

	// Auth middleware (for protected routes)
	authMiddleware := middleware.AuthMiddleware(authSvc)

	// API routes
	api := r.Group("/api")

	// Auth routes (public - no auth required)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)

	// Protected routes - all require JWT authentication
	protected := api.Group("")
	protected.Use(authMiddleware)
	{
		// Auth
		protected.GET("/auth/me", authHandler.Me)

		// Category routes
		categoryRepo := repository.NewCategoryRepository(db)
		categorySvc := service.NewCategoryService(categoryRepo)
		categoryHandler := handler.NewCategoryHandler(categorySvc)
		protected.GET("/categories", categoryHandler.List)
		protected.GET("/categories/:id", categoryHandler.Get)
		protected.POST("/categories", categoryHandler.Create)
		protected.PUT("/categories/:id", categoryHandler.Update)
		protected.DELETE("/categories/:id", categoryHandler.Delete)

		// Staff routes
		staffSvc := service.NewStaffService(staffRepo)
		staffHandler := handler.NewStaffHandler(staffSvc)
		protected.GET("/staff", staffHandler.List)
		protected.GET("/staff/:id", staffHandler.Get)
		protected.POST("/staff", staffHandler.Create)
		protected.PUT("/staff/:id", staffHandler.Update)
		protected.DELETE("/staff/:id", staffHandler.Delete)
		protected.PUT("/staff/:id/status", staffHandler.UpdateStatus)

		// User routes
		protected.GET("/users", userHandler.List)
		protected.GET("/users/:id", userHandler.Get)
		protected.POST("/users", userHandler.Create)
		protected.PUT("/users/:id", userHandler.Update)
		protected.DELETE("/users/:id", userHandler.Delete)

		// Product routes
		protected.GET("/products", productHandler.List)
		protected.GET("/products/:id", productHandler.Get)
		protected.POST("/products", productHandler.Create)
		protected.PUT("/products/:id", productHandler.Update)
		protected.DELETE("/products/:id", productHandler.Delete)

		// Order routes
		protected.GET("/orders", orderHandler.List)
		protected.GET("/orders/:id", orderHandler.Get)
		protected.POST("/orders", orderHandler.Create)
		protected.DELETE("/orders/:id", orderHandler.Delete)
		protected.PUT("/orders/:id/staff", orderHandler.AssignStaff)
		protected.POST("/orders/:id/paid", orderHandler.Paid)
		protected.POST("/orders/:id/ship", orderHandler.Ship)
		protected.POST("/orders/:id/complete", orderHandler.Complete)
		protected.POST("/orders/:id/cancel", orderHandler.Cancel)

		// Stats routes
		protected.GET("/stats/orders", statsHandler.OrderStats)
		protected.GET("/stats/revenue", statsHandler.RevenueStats)
	}
}
