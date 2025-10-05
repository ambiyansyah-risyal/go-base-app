package database

import (
	"database/sql"
	"fmt"

	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
)

// Migrator handles database migrations
type Migrator struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, log *logger.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: log,
	}
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	m.logger.Info("Running database migrations...")

	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migrations to run
	migrations := m.getMigrations()

	for _, migration := range migrations {
		// Check if migration has already been run
		if exists, err := m.migrationExists(migration.Version); err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		} else if exists {
			m.logger.Debug("Migration already applied", "version", migration.Version)
			continue
		}

		// Run the migration
		m.logger.Info("Applying migration", "version", migration.Version, "description", migration.Description)
		
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute migration SQL
		if _, err := tx.Exec(migration.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}

		// Record migration
		if err := m.recordMigration(tx, migration.Version, migration.Description); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration: %w", err)
		}

		m.logger.Info("Migration applied successfully", "version", migration.Version)
	}

	m.logger.Info("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(255) PRIMARY KEY,
			description TEXT,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	if _, err := m.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// migrationExists checks if a migration has been applied
func (m *Migrator) migrationExists(version string) (bool, error) {
	query := "SELECT COUNT(*) FROM migrations WHERE version = ?"
	
	var count int
	if err := m.db.QueryRow(query, version).Scan(&count); err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// recordMigration records that a migration has been applied
func (m *Migrator) recordMigration(tx *sql.Tx, version, description string) error {
	query := "INSERT INTO migrations (version, description) VALUES (?, ?)"
	
	_, err := tx.Exec(query, version, description)
	return err
}

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	SQL         string
}

// getMigrations returns all available migrations
func (m *Migrator) getMigrations() []Migration {
	return []Migration{
		{
			Version:     "001_create_users_table",
			Description: "Create users table",
			SQL: `
				CREATE TABLE users (
					id TEXT PRIMARY KEY,
					email TEXT UNIQUE NOT NULL,
					username TEXT UNIQUE NOT NULL,
					password TEXT NOT NULL,
					first_name TEXT NOT NULL,
					last_name TEXT NOT NULL,
					is_active BOOLEAN DEFAULT TRUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP NULL
				)
			`,
		},
		{
			Version:     "002_create_users_indexes",
			Description: "Create indexes on users table",
			SQL: `
				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
				CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
				CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
			`,
		},
	}
}