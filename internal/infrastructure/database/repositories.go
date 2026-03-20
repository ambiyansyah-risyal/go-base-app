package database

import (
	"database/sql"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/repository"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
)

// Repositories holds all repository implementations
type Repositories struct {
	User   repository.UserRepository
	Health repository.HealthRepository
}

// NewRepositories creates a new repositories instance
func NewRepositories(db *sql.DB, log *logger.Logger) *Repositories {
	return &Repositories{
		User:   NewUserRepository(db, log),
		Health: NewHealthRepository(db, log),
	}
}
