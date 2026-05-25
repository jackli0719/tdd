package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port string
	DSN  string
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
	return &Config{
		Port: port,
		DSN:  dsn,
	}
}
