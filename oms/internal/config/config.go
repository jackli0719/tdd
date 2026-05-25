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
		dsn = "root:password@tcp(localhost:3306)/oms?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return &Config{
		Port: port,
		DSN:  dsn,
	}
}
