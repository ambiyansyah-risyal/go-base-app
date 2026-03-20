package usecase

import (
	"github.com/ambiyansyah-risyal/go-base-app/internal/infrastructure/database"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
)

// Usecases holds all use cases
type Usecases struct {
	Health HealthUsecase
	User   UserUsecase
	Auth   AuthUsecase
}

// NewUsecases creates a new usecases instance
func NewUsecases(repos *database.Repositories, log *logger.Logger, validator *validator.Validator) *Usecases {
	return &Usecases{
		Health: NewHealthUsecase(repos.Health, log),
		User:   NewUserUsecase(repos.User, log, validator),
		Auth:   NewAuthUsecase(repos.User, log, validator),
	}
}
