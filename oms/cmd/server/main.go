package main

import (
	"log"

	"oms/internal/config"
	"oms/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database (optional - continue without DB for now)
	db, err := config.InitDB(cfg.DSN)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Server will start without database connection")
	}

	// Setup Gin
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Setup routes
	router.Setup(r, db)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
