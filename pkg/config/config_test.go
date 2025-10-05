package config

import "testing"

func TestLoad(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}
	if cfg.App.Name == "" {
		t.Error("App name should not be empty")
	}
}

func TestConfigValidation(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Name:        "test-app",
			Version:     "1.0.0",
			Environment: "test",
			Debug:       false,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    ":memory:",
		},
		Logger: LoggerConfig{
			Level: "info",
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}

func TestInvalidPortValidation(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Name: "test",
		},
		Server: ServerConfig{
			Port: -1, // Invalid port
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    ":memory:",
		},
		Logger: LoggerConfig{
			Level: "info",
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Error("Expected validation to fail for invalid port")
	}
}
