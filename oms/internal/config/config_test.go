package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Clear env
	os.Unsetenv("PORT")
	os.Unsetenv("DSN")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}

	if cfg.DSN == "" {
		t.Error("expected non-empty DSN")
	}
}

func TestLoadWithEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DSN", "test_dsn")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DSN")
	}()

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.DSN != "test_dsn" {
		t.Errorf("expected DSN test_dsn, got %s", cfg.DSN)
	}
}
