package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig      `json:"app" yaml:"app"`
	Server   ServerConfig   `json:"server" yaml:"server"`
	Database DatabaseConfig `json:"database" yaml:"database"`
	Security SecurityConfig `json:"security" yaml:"security"`
	Logger   LoggerConfig   `json:"logger" yaml:"logger"`
	Features FeatureConfig  `json:"features" yaml:"features"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Environment string `json:"environment" yaml:"environment"`
	Debug       bool   `json:"debug" yaml:"debug"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	GracefulStop time.Duration `json:"graceful_stop" yaml:"graceful_stop"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver          string        `json:"driver" yaml:"driver"`
	DSN             string        `json:"dsn" yaml:"dsn"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	AutoMigrate     bool          `json:"auto_migrate" yaml:"auto_migrate"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret         string        `json:"jwt_secret" yaml:"jwt_secret"`
	JWTExpiration     time.Duration `json:"jwt_expiration" yaml:"jwt_expiration"`
	RateLimitRequests int           `json:"rate_limit_requests" yaml:"rate_limit_requests"`
	RateLimitWindow   time.Duration `json:"rate_limit_window" yaml:"rate_limit_window"`
	CORSAllowOrigins  []string      `json:"cors_allow_origins" yaml:"cors_allow_origins"`
	EnableHTTPS       bool          `json:"enable_https" yaml:"enable_https"`
	TLSCertFile       string        `json:"tls_cert_file" yaml:"tls_cert_file"`
	TLSKeyFile        string        `json:"tls_key_file" yaml:"tls_key_file"`
}

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"`
	Output     string `json:"output" yaml:"output"`
	Structured bool   `json:"structured" yaml:"structured"`
}

// FeatureConfig holds feature flag configuration
type FeatureConfig struct {
	EnableMetrics    bool `json:"enable_metrics" yaml:"enable_metrics"`
	EnableTracing    bool `json:"enable_tracing" yaml:"enable_tracing"`
	EnableProfiling  bool `json:"enable_profiling" yaml:"enable_profiling"`
	EnableSwagger    bool `json:"enable_swagger" yaml:"enable_swagger"`
	EnableHealthz    bool `json:"enable_healthz" yaml:"enable_healthz"`
	EnablePlayground bool `json:"enable_playground" yaml:"enable_playground"`
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:        getEnvString("APP_NAME", "go-base-app"),
			Version:     getEnvString("APP_VERSION", "dev"),
			Environment: getEnvString("APP_ENVIRONMENT", "development"),
			Debug:       getEnvBool("APP_DEBUG", true),
		},
		Server: ServerConfig{
			Host:         getEnvString("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			GracefulStop: getEnvDuration("SERVER_GRACEFUL_STOP", 10*time.Second),
		},
		Database: DatabaseConfig{
			Driver:          getEnvString("DB_DRIVER", "sqlite"),
			DSN:             getEnvString("DB_DSN", "./data/app.db"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
			AutoMigrate:     getEnvBool("DB_AUTO_MIGRATE", true),
		},
		Security: SecurityConfig{
			JWTSecret:         getEnvString("JWT_SECRET", generateDefaultSecret()),
			JWTExpiration:     getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvDuration("RATE_LIMIT_WINDOW", time.Minute),
			CORSAllowOrigins:  getEnvStringSlice("CORS_ALLOW_ORIGINS", []string{"*"}),
			EnableHTTPS:       getEnvBool("ENABLE_HTTPS", false),
			TLSCertFile:       getEnvString("TLS_CERT_FILE", ""),
			TLSKeyFile:        getEnvString("TLS_KEY_FILE", ""),
		},
		Logger: LoggerConfig{
			Level:      getEnvString("LOG_LEVEL", "info"),
			Format:     getEnvString("LOG_FORMAT", "json"),
			Output:     getEnvString("LOG_OUTPUT", "stdout"),
			Structured: getEnvBool("LOG_STRUCTURED", true),
		},
		Features: FeatureConfig{
			EnableMetrics:    getEnvBool("FEATURE_METRICS", true),
			EnableTracing:    getEnvBool("FEATURE_TRACING", false),
			EnableProfiling:  getEnvBool("FEATURE_PROFILING", false),
			EnableSwagger:    getEnvBool("FEATURE_SWAGGER", true),
			EnableHealthz:    getEnvBool("FEATURE_HEALTHZ", true),
			EnablePlayground: getEnvBool("FEATURE_PLAYGROUND", false),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if c.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[strings.ToLower(c.Logger.Level)] {
		return fmt.Errorf("invalid log level: %s", c.Logger.Level)
	}

	return nil
}

// GetDSN returns the formatted database connection string
func (c *Config) GetDSN() string {
	return c.Database.DSN
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.App.Environment) == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.App.Environment) == "development"
}

// Helper functions for environment variable parsing
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func generateDefaultSecret() string {
	return "default-secret-change-in-production"
}