package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"    // SQLite driver
)

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig, log *logger.Logger) (*sql.DB, error) {
	var driverName string
	var dataSourceName string

	switch cfg.Driver {
	case "sqlite", "sqlite3":
		driverName = "sqlite3"
		dataSourceName = cfg.DSN
	case "postgres", "postgresql":
		driverName = "postgres"
		dataSourceName = cfg.DSN
	case "mysql":
		driverName = "mysql"
		dataSourceName = cfg.DSN
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established",
		"driver", cfg.Driver,
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
	)

	return db, nil
}

// GetDatabaseStats returns database connection statistics
func GetDatabaseStats(db *sql.DB) map[string]interface{} {
	stats := db.Stats()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// HealthCheck performs a simple health check on the database
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
