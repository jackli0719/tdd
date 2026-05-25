package router

import (
	"oms/internal/handler"
	"oms/internal/repository"
	"oms/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup configures all routes
func Setup(r *gin.Engine, db *gorm.DB) {
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
		// User routes
		userRepo := repository.NewUserRepository(db)
		userSvc := service.NewUserService(userRepo)
		userHandler := handler.NewUserHandler(userSvc)
		api.GET("/users", userHandler.List)
		api.GET("/users/:id", userHandler.Get)
		api.POST("/users", userHandler.Create)
		api.PUT("/users/:id", userHandler.Update)
		api.DELETE("/users/:id", userHandler.Delete)
	}
}
