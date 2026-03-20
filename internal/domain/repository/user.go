package repository

import (
	"context"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with filtering and pagination
	List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, int64, error)

	// ExistsByEmail checks if a user exists by email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

// HealthRepository defines the interface for health checks
type HealthRepository interface {
	// Ping checks if the database is reachable
	Ping(ctx context.Context) error

	// GetStats returns database statistics
	GetStats(ctx context.Context) (map[string]interface{}, error)
}
