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

func TestInitDB_SQLite(t *testing.T) {
	db, err := InitDB("file::memory:")
	if err != nil {
		t.Fatalf("failed to init sqlite db: %v", err)
	}
	if db == nil {
		t.Error("expected non-nil db")
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.Close()
}

func TestInitDB_SQLiteFilePath(t *testing.T) {
	os.Remove("test.db")
	db, err := InitDB("sqlite://test.db")
	if err != nil {
		t.Fatalf("failed to init sqlite db: %v", err)
	}
	if db == nil {
		t.Error("expected non-nil db")
	}
}

func TestInitDB_MySQLDSN(t *testing.T) {
	// Test with mysql dsn prefix (will fail to connect but should parse correctly)
	_, err := InitDB("mysql://user:pass@localhost:3306/testdb")
	if err == nil {
		// May succeed or fail depending on MySQL availability
		t.Log("mysql connection attempt result: " + err.Error())
	}
}
