package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance based on configuration
func New(cfg config.LoggerConfig) *Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var output io.Writer
	switch strings.ToLower(cfg.Output) {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output = os.Stdout
	}

	var handler slog.Handler
	if cfg.Structured {
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := slog.New(handler)
	return &Logger{Logger: logger}
}

// WithContext adds contextual fields to the logger
func (l *Logger) WithContext(args ...any) *Logger {
	return &Logger{Logger: l.With(args...)}
}

// WithGroup adds a group to the logger
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{Logger: l.Logger.WithGroup(name)}
}

// HTTP creates a logger instance for HTTP requests
func (l *Logger) HTTP() *Logger {
	return l.WithGroup("http")
}

// Database creates a logger instance for database operations
func (l *Logger) Database() *Logger {
	return l.WithGroup("database")
}

// Auth creates a logger instance for authentication operations
func (l *Logger) Auth() *Logger {
	return l.WithGroup("auth")
}

// Service creates a logger instance for service operations
func (l *Logger) Service() *Logger {
	return l.WithGroup("service")
}

// Infrastructure creates a logger instance for infrastructure operations
func (l *Logger) Infrastructure() *Logger {
	return l.WithGroup("infrastructure")
}

// Request logs an HTTP request
func (l *Logger) Request(method, path, userAgent, clientIP string, statusCode int, latency string) {
	l.HTTP().Info("request",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"user_agent", userAgent,
		"client_ip", clientIP,
		"latency", latency,
	)
}

// QueryLog logs a database query
func (l *Logger) QueryLog(query string, duration string, err error) {
	if err != nil {
		l.Database().Error("query failed",
			"query", query,
			"duration", duration,
			"error", err,
		)
	} else {
		l.Database().Debug("query executed",
			"query", query,
			"duration", duration,
		)
	}
}

// AuthLog logs authentication events
func (l *Logger) AuthLog(action, userID, clientIP string, success bool, reason string) {
	logger := l.Auth().With(
		"action", action,
		"user_id", userID,
		"client_ip", clientIP,
		"success", success,
	)

	if success {
		logger.Info("authentication successful")
	} else {
		logger.Warn("authentication failed", "reason", reason)
	}
}

// StartupLog logs application startup information
func (l *Logger) StartupLog(appName, version, environment string, port int, features map[string]bool) {
	l.Info("application starting",
		"app_name", appName,
		"version", version,
		"environment", environment,
		"port", port,
		"features", features,
	)
}

// ShutdownLog logs application shutdown information
func (l *Logger) ShutdownLog(appName string, reason string) {
	l.Info("application shutting down",
		"app_name", appName,
		"reason", reason,
	)
}

// ErrorWithContext logs an error with additional context
func (l *Logger) ErrorWithContext(msg string, err error, context map[string]any) {
	args := []any{"error", err}
	for k, v := range context {
		args = append(args, k, v)
	}
	l.Error(msg, args...)
}

// Default logger instance
var defaultLogger *Logger

// Init initializes the default logger
func Init(cfg config.LoggerConfig) {
	defaultLogger = New(cfg)
}

// Default returns the default logger instance
func Default() *Logger {
	if defaultLogger == nil {
		// Fallback to a basic logger if not initialized
		defaultLogger = New(config.LoggerConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			Structured: true,
		})
	}
	return defaultLogger
}

// Package-level convenience functions
func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func With(args ...any) *Logger {
	return Default().WithContext(args...)
}

func WithGroup(name string) *Logger {
	return Default().WithGroup(name)
}