package logger

import (
	"testing"

	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
)

func TestInit(t *testing.T) {
	cfg := config.LoggerConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		Structured: true,
	}

	Init(cfg)

	logger := Default()
	if logger == nil {
		t.Error("Default logger should not be nil after Init")
	}
}

func TestNew(t *testing.T) {
	cfg := config.LoggerConfig{
		Level:      "debug",
		Format:     "text",
		Output:     "stdout",
		Structured: false,
	}

	logger := New(cfg)
	if logger == nil {
		t.Error("New should not return nil")
	}
}

func TestWithGroup(t *testing.T) {
	cfg := config.LoggerConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		Structured: true,
	}

	logger := New(cfg)
	grouped := logger.WithGroup("test")

	if grouped == nil {
		t.Error("WithGroup should not return nil")
	}
}

func TestLoggerMethods(t *testing.T) {
	cfg := config.LoggerConfig{
		Level:      "debug",
		Format:     "json",
		Output:     "stdout",
		Structured: true,
	}

	Init(cfg)

	// Test that these don't panic
	Info("test info message", "key", "value")
	Debug("test debug message", "key", "value")
	Warn("test warn message", "key", "value")
	Error("test error message", "key", "value")
}

func TestSpecializedLoggers(t *testing.T) {
	cfg := config.LoggerConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		Structured: true,
	}

	Init(cfg)

	// Test that we can create grouped loggers
	logger := Default()
	httpLogger := logger.WithGroup("http")
	if httpLogger == nil {
		t.Error("HTTP logger should not be nil")
	}

	dbLogger := logger.WithGroup("database")
	if dbLogger == nil {
		t.Error("Database logger should not be nil")
	}

	authLogger := logger.WithGroup("auth")
	if authLogger == nil {
		t.Error("Auth logger should not be nil")
	}
}