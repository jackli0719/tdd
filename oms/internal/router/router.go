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

	// Services
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, jwtSecret)
	productSvc := service.NewProductService(productRepo)
	orderSvc := service.NewOrderService(orderRepo, userRepo, productRepo)

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
	{
		// Auth routes (public)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)

		// Auth routes (protected)
		auth := api.Group("/auth")
		auth.Use(authMiddleware)
		{
			auth.GET("/me", authHandler.Me)
		}

		// Category routes (public for now - can be protected later)
		categoryRepo := repository.NewCategoryRepository(db)
		categorySvc := service.NewCategoryService(categoryRepo)
		categoryHandler := handler.NewCategoryHandler(categorySvc)
		api.GET("/categories", categoryHandler.List)
		api.GET("/categories/:id", categoryHandler.Get)
		api.POST("/categories", categoryHandler.Create)
		api.PUT("/categories/:id", categoryHandler.Update)
		api.DELETE("/categories/:id", categoryHandler.Delete)

		// User routes (protected)
		users := api.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		// Product routes (public for now)
		api.GET("/products", productHandler.List)
		api.GET("/products/:id", productHandler.Get)
		api.POST("/products", productHandler.Create)
		api.PUT("/products/:id", productHandler.Update)
		api.DELETE("/products/:id", productHandler.Delete)

		// Order routes (public for now)
		api.GET("/orders", orderHandler.List)
		api.GET("/orders/:id", orderHandler.Get)
		api.POST("/orders", orderHandler.Create)
		api.DELETE("/orders/:id", orderHandler.Delete)
		api.POST("/orders/:id/confirm", orderHandler.Paid)
		api.POST("/orders/:id/start", orderHandler.Ship)
		api.POST("/orders/:id/complete", orderHandler.Complete)
		api.POST("/orders/:id/cancel", orderHandler.Cancel)

		// Stats routes (public for now)
		api.GET("/stats/orders", statsHandler.OrderStats)
		api.GET("/stats/revenue", statsHandler.RevenueStats)
	}
}
