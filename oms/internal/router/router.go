package router

import (
	"github.com/gin-gonic/gin"
)

// Setup configures all routes
func Setup(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// TODO: Register user, product, order routes
		api.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "users placeholder"})
		})
	}
}
