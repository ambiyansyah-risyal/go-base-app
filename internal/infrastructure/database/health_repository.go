package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/repository"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
)

// healthRepository implements repository.HealthRepository
type healthRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewHealthRepository creates a new health repository
func NewHealthRepository(db *sql.DB, log *logger.Logger) repository.HealthRepository {
	return &healthRepository{
		db:     db,
		logger: log,
	}
}

// Ping checks if the database is reachable
func (r *healthRepository) Ping(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		r.logger.Database().Error("Database ping failed", "error", err)
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}

// GetStats returns database statistics
func (r *healthRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := GetDatabaseStats(r.db)

	// Add additional database-specific statistics
	stats["ping_ok"] = true

	// Try to get database version or other metadata
	var version string
	if err := r.db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version); err == nil {
		stats["version"] = version
		stats["engine"] = "sqlite"
	} else {
		// Try PostgreSQL version
		if err := r.db.QueryRowContext(ctx, "SELECT version()").Scan(&version); err == nil {
			stats["version"] = version
			stats["engine"] = "postgresql"
		} else {
			// Try MySQL version
			if err := r.db.QueryRowContext(ctx, "SELECT @@version").Scan(&version); err == nil {
				stats["version"] = version
				stats["engine"] = "mysql"
			}
		}
	}

	return stats, nil
}
