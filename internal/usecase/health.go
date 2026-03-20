package usecase

import (
	"context"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/repository"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
)

// healthUsecase implements HealthUsecase interface
type healthUsecase struct {
	healthRepo repository.HealthRepository
	logger     *logger.Logger
}

// NewHealthUsecase creates a new health usecase
func NewHealthUsecase(healthRepo repository.HealthRepository, log *logger.Logger) HealthUsecase {
	return &healthUsecase{
		healthRepo: healthRepo,
		logger:     log,
	}
}

// GetHealth returns the application health status
func (u *healthUsecase) GetHealth(ctx context.Context) map[string]interface{} {
	health := map[string]interface{}{
		"status": "ok",
		"checks": map[string]interface{}{},
	}

	// Check database health
	if err := u.healthRepo.Ping(ctx); err != nil {
		health["status"] = "degraded"
		health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	} else {
		health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
			"status": "healthy",
		}
	}

	// Add database statistics if available
	if stats, err := u.healthRepo.GetStats(ctx); err == nil {
		health["database_stats"] = stats
	}

	return health
}

// IsLive returns true if the application is alive
func (u *healthUsecase) IsLive(ctx context.Context) bool {
	// Basic liveness check - if we can execute this function, we're alive
	return true
}

// IsReady returns true if the application is ready to serve requests
func (u *healthUsecase) IsReady(ctx context.Context) bool {
	// Check if database is accessible
	if err := u.healthRepo.Ping(ctx); err != nil {
		u.logger.Warn("Readiness check failed", "error", err)
		return false
	}

	return true
}
