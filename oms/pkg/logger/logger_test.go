package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug level", "debug"},
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"unknown level defaults to info", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.level)
			if err != nil {
				t.Errorf("Init() error = %v", err)
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Initialize first
	Init("info")
	logger := GetLogger()
	if logger == nil {
		t.Error("GetLogger() returned nil")
	}
}

func TestSync(t *testing.T) {
	// Should not panic even if logger is nil
	Sync()
}

func TestDebug(t *testing.T) {
	Init("debug")
	Debug("test debug message", zap.String("key", "value"))
}

func TestInfo(t *testing.T) {
	Init("info")
	Info("test info message", zap.String("key", "value"))
}

func TestWarn(t *testing.T) {
	Init("warn")
	Warn("test warn message", zap.String("key", "value"))
}

func TestError(t *testing.T) {
	Init("error")
	Error("test error message", zap.String("key", "value"))
}

func TestFatal(t *testing.T) {
	// Fatal calls os.Exit, so we need to test it differently
	t.Run("fatal does not return", func(t *testing.T) {
		done := make(chan bool, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// Expected - Fatal calls os.Exit
				}
				done <- true
			}()
			Fatal("test fatal message", zap.String("key", "value"))
			done <- true
		}()
		<-done
	})
}

func TestWith(t *testing.T) {
	Init("info")
	logger := With(zap.String("key", "value"))
	if logger == nil {
		t.Error("With() returned nil")
	}
}

func TestGetLoggerDefault(t *testing.T) {
	// Reset log to nil
	log = nil
	logger := GetLogger()
	if logger == nil {
		t.Error("GetLogger() returned nil after reset")
	}
}