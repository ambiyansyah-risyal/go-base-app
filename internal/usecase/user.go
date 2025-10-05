package usecase

import (
	"context"
	"fmt"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/repository"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo  repository.UserRepository
	logger    *logger.Logger
	validator *validator.Validator
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository, log *logger.Logger, validator *validator.Validator) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		logger:    log,
		validator: validator,
	}
}

// GetByID retrieves a user by ID
func (u *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Failed to get user by ID", "error", err, "user_id", id)
		return nil, err
	}

	if user.IsDeleted() {
		return nil, entity.ErrUserDeleted
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.IsDeleted() {
		return nil, entity.ErrUserDeleted
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (u *userUsecase) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user.IsDeleted() {
		return nil, entity.ErrUserDeleted
	}

	return user, nil
}

// Create creates a new user
func (u *userUsecase) Create(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	// Validate request
	if err := u.validator.Validate(&req); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Failed to check email existence", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to validate email: %w", err)
	}
	if exists {
		return nil, entity.ErrUserEmailExists
	}

	// Check if username already exists
	exists, err = u.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		u.logger.Error("Failed to check username existence", "error", err, "username", req.Username)
		return nil, fmt.Errorf("failed to validate username: %w", err)
	}
	if exists {
		return nil, entity.ErrUserUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := entity.NewUser(req.Email, req.Username, string(hashedPassword), req.FirstName, req.LastName)

	// Save to repository
	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.Error("Failed to create user", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	u.logger.Info("User created successfully", "user_id", user.ID, "email", req.Email)
	return user, nil
}

// Update updates an existing user
func (u *userUsecase) Update(ctx context.Context, id uuid.UUID, req entity.UpdateUserRequest) (*entity.User, error) {
	// Validate request
	if err := u.validator.Validate(&req); err != nil {
		return nil, err
	}

	// Get existing user
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.IsDeleted() {
		return nil, entity.ErrUserDeleted
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.IsActive != nil {
		if *req.IsActive {
			user.Activate()
		} else {
			user.Deactivate()
		}
	}

	// Save changes
	if err := u.userRepo.Update(ctx, user); err != nil {
		u.logger.Error("Failed to update user", "error", err, "user_id", id)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	u.logger.Info("User updated successfully", "user_id", id)
	return user, nil
}

// Delete soft deletes a user
func (u *userUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user.IsDeleted() {
		return entity.ErrUserDeleted
	}

	// Soft delete
	if err := u.userRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Failed to delete user", "error", err, "user_id", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	u.logger.Info("User deleted successfully", "user_id", id)
	return nil
}

// List retrieves users with filtering and pagination
func (u *userUsecase) List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, int64, error) {
	users, total, err := u.userRepo.List(ctx, filter)
	if err != nil {
		u.logger.Error("Failed to list users", "error", err, "filter", filter)
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
