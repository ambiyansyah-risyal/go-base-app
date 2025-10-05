package handler

import (
	"github.com/ambiyansyah-risyal/go-base-app/internal/usecase"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Health *HealthHandler
	User   *UserHandler
	Auth   *AuthHandler
}

// NewHandlers creates a new handlers instance
func NewHandlers(usecases *usecase.Usecases, log *logger.Logger, validator *validator.Validator) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(usecases.Health, log),
		User:   NewUserHandler(usecases.User, log, validator),
		Auth:   NewAuthHandler(usecases.Auth, log, validator),
	}
}