package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/repository"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// authUsecase implements AuthUsecase interface
type authUsecase struct {
	userRepo  repository.UserRepository
	logger    *logger.Logger
	validator *validator.Validator
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo repository.UserRepository, log *logger.Logger, validator *validator.Validator) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		logger:    log,
		validator: validator,
	}
}

// Login authenticates a user and returns tokens
func (a *authUsecase) Login(ctx context.Context, email, password string) (*entity.User, *AuthTokens, error) {
	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if entity.IsNotFoundError(err) {
			return nil, nil, entity.ErrInvalidCredentials
		}
		a.logger.Error("Failed to get user by email during login", "error", err, "email", email)
		return nil, nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	// Check if user is deleted
	if user.IsDeleted() {
		return nil, nil, entity.ErrUserDeleted
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, entity.ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, entity.ErrInvalidCredentials
	}

	// Generate tokens
	tokens, err := a.generateTokens(user)
	if err != nil {
		a.logger.Error("Failed to generate tokens", "error", err, "user_id", user.ID)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// Register creates a new user account and returns tokens
func (a *authUsecase) Register(ctx context.Context, req entity.CreateUserRequest) (*entity.User, *AuthTokens, error) {
	// Validate request
	if err := a.validator.Validate(&req); err != nil {
		return nil, nil, err
	}

	// Check if email already exists
	exists, err := a.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		a.logger.Error("Failed to check email existence during registration", "error", err, "email", req.Email)
		return nil, nil, fmt.Errorf("failed to validate email: %w", err)
	}
	if exists {
		return nil, nil, entity.ErrUserEmailExists
	}

	// Check if username already exists
	exists, err = a.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		a.logger.Error("Failed to check username existence during registration", "error", err, "username", req.Username)
		return nil, nil, fmt.Errorf("failed to validate username: %w", err)
	}
	if exists {
		return nil, nil, entity.ErrUserUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("Failed to hash password during registration", "error", err)
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := entity.NewUser(req.Email, req.Username, string(hashedPassword), req.FirstName, req.LastName)

	// Save to repository
	if err := a.userRepo.Create(ctx, user); err != nil {
		a.logger.Error("Failed to create user during registration", "error", err, "email", req.Email)
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := a.generateTokens(user)
	if err != nil {
		a.logger.Error("Failed to generate tokens during registration", "error", err, "user_id", user.ID)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	a.logger.Info("User registered successfully", "user_id", user.ID, "email", req.Email)
	return user, tokens, nil
}

// RefreshTokens refreshes authentication tokens
func (a *authUsecase) RefreshTokens(ctx context.Context, refreshToken string) (*entity.User, *AuthTokens, error) {
	// In a full implementation, you would:
	// 1. Validate the refresh token (decode JWT, check signature, expiration)
	// 2. Extract user ID from the token
	// 3. Get user from database
	// 4. Generate new tokens
	// 5. Optionally invalidate the old refresh token

	// For now, return an error indicating not implemented
	return nil, nil, entity.ErrNotImplemented
}

// ValidateToken validates an access token and returns the user
func (a *authUsecase) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	// In a full implementation, you would:
	// 1. Parse and validate JWT token
	// 2. Extract user ID from token claims
	// 3. Get user from database
	// 4. Check if user is still active

	// For now, return an error indicating not implemented
	return nil, entity.ErrNotImplemented
}

// Logout invalidates user tokens
func (a *authUsecase) Logout(ctx context.Context, userID uuid.UUID) error {
	// In a full implementation, you would:
	// 1. Add tokens to a blacklist
	// 2. Or remove refresh tokens from database
	// 3. Log the logout event

	a.logger.Info("User logged out", "user_id", userID)
	return nil
}

// generateTokens generates access and refresh tokens for a user
func (a *authUsecase) generateTokens(user *entity.User) (*AuthTokens, error) {
	// In a full implementation, you would:
	// 1. Create JWT access token with user claims and short expiration
	// 2. Create JWT refresh token with longer expiration
	// 3. Sign tokens with secret key
	// 4. Store refresh token in database (optional)

	// For demonstration purposes, return mock tokens
	now := time.Now()
	expiresIn := int64(24 * time.Hour / time.Second) // 24 hours in seconds

	tokens := &AuthTokens{
		AccessToken:  fmt.Sprintf("mock_access_token_%s_%d", user.ID.String(), now.Unix()),
		RefreshToken: fmt.Sprintf("mock_refresh_token_%s_%d", user.ID.String(), now.Unix()),
		ExpiresIn:    expiresIn,
	}

	return tokens, nil
}
