package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/repository"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/google/uuid"
)

// userRepository implements repository.UserRepository
type userRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB, log *logger.Logger) repository.UserRepository {
	return &userRepository{
		db:     db,
		logger: log,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, username, password, first_name, last_name, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID.String(),
		user.Email,
		user.Username,
		user.Password,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		r.logger.Database().Error("Failed to create user", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, username, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	user := &entity.User{}
	var deletedAt sql.NullTime
	var userIDStr string

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&userIDStr,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		r.logger.Database().Error("Failed to get user by ID", "error", err, "user_id", id)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse UUID from string
	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.logger.Database().Error("Failed to parse user ID", "error", err, "user_id_str", userIDStr)
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}
	user.ID = parsedID

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, username, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	user := &entity.User{}
	var deletedAt sql.NullTime
	var userIDStr string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&userIDStr,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		r.logger.Database().Error("Failed to get user by email", "error", err, "email", email)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse UUID from string
	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.logger.Database().Error("Failed to parse user ID", "error", err, "user_id_str", userIDStr)
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}
	user.ID = parsedID

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, email, username, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE username = ? AND deleted_at IS NULL
	`

	user := &entity.User{}
	var deletedAt sql.NullTime
	var userIDStr string

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&userIDStr,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		r.logger.Database().Error("Failed to get user by username", "error", err, "username", username)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse UUID from string
	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.logger.Database().Error("Failed to parse user ID", "error", err, "user_id_str", userIDStr)
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}
	user.ID = parsedID

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = ?, username = ?, password = ?, first_name = ?, last_name = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.ID.String(),
	)

	if err != nil {
		r.logger.Database().Error("Failed to update user", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = ?, is_active = false, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, now, id.String())

	if err != nil {
		r.logger.Database().Error("Failed to delete user", "error", err, "user_id", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// List retrieves users with filtering and pagination
func (r *userRepository) List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if filter.Email != "" {
		whereClause += fmt.Sprintf(" AND email LIKE ?")
		args = append(args, "%"+filter.Email+"%")
		argIndex++
	}

	if filter.Username != "" {
		whereClause += fmt.Sprintf(" AND username LIKE ?")
		args = append(args, "%"+filter.Username+"%")
		argIndex++
	}

	if filter.IsActive != nil {
		whereClause += fmt.Sprintf(" AND is_active = ?")
		args = append(args, *filter.IsActive)
		argIndex++
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM users " + whereClause
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Database().Error("Failed to count users", "error", err)
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "ORDER BY " + filter.SortBy
	if filter.SortDesc {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	// Build main query
	query := `
		SELECT id, email, username, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		` + whereClause + " " + orderBy + " LIMIT ? OFFSET ?"

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Database().Error("Failed to list users", "error", err, "filter", filter)
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		var deletedAt sql.NullTime
		var userIDStr string

		err := rows.Scan(
			&userIDStr,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&deletedAt,
		)

		if err != nil {
			r.logger.Database().Error("Failed to scan user row", "error", err)
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		// Parse UUID from string
		parsedID, err := uuid.Parse(userIDStr)
		if err != nil {
			r.logger.Database().Error("Failed to parse user ID", "error", err, "user_id_str", userIDStr)
			return nil, 0, fmt.Errorf("failed to parse user ID: %w", err)
		}
		user.ID = parsedID

		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		r.logger.Database().Error("Error iterating user rows", "error", err)
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		r.logger.Database().Error("Failed to check user existence by email", "error", err, "email", email)
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// ExistsByUsername checks if a user exists by username
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND deleted_at IS NULL)"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		r.logger.Database().Error("Failed to check user existence by username", "error", err, "username", username)
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
