package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port      string
	DSN       string
	JWTSecret string
}

// Load reads configuration from environment variables
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DSN")
	if dsn == "" {
		// Default to SQLite for local development
		dsn = "sqlite://oms.db"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}
	return &Config{
		Port:      port,
		DSN:       dsn,
		JWTSecret: jwtSecret,
	}
}
