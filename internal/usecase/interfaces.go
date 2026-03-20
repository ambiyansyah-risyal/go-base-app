package usecase

import (
	"context"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/google/uuid"
)

// AuthTokens represents authentication tokens
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// HealthUsecase defines health check operations
type HealthUsecase interface {
	GetHealth(ctx context.Context) map[string]interface{}
	IsLive(ctx context.Context) bool
	IsReady(ctx context.Context) bool
}

// UserUsecase defines user-related operations
type UserUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, req entity.UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, int64, error)
}

// AuthUsecase defines authentication operations
type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*entity.User, *AuthTokens, error)
	Register(ctx context.Context, req entity.CreateUserRequest) (*entity.User, *AuthTokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*entity.User, *AuthTokens, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}
