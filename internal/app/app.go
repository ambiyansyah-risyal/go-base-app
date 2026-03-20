package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ambiyansyah-risyal/go-base-app/internal/infrastructure/database"
	"github.com/ambiyansyah-risyal/go-base-app/internal/infrastructure/http/handler"
	"github.com/ambiyansyah-risyal/go-base-app/internal/infrastructure/http/router"
	"github.com/ambiyansyah-risyal/go-base-app/internal/usecase"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/middleware"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
	"github.com/gin-gonic/gin"
)

// App represents the application
type App struct {
	config    *config.Config
	logger    *logger.Logger
	db        *sql.DB
	router    *gin.Engine
	validator *validator.Validator
	usecases  *usecase.Usecases
	handlers  *handler.Handlers
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, log *logger.Logger) (*App, error) {
	app := &App{
		config: cfg,
		logger: log,
	}

	// Initialize validator
	app.validator = validator.New()

	// Initialize database
	if err := app.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	repositories, err := app.initRepositories()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	// Initialize use cases
	app.usecases = usecase.NewUsecases(repositories, app.logger, app.validator)

	// Initialize handlers
	app.handlers = handler.NewHandlers(app.usecases, app.logger, app.validator)

	// Initialize router
	if err := app.initRouter(); err != nil {
		return nil, fmt.Errorf("failed to initialize router: %w", err)
	}

	return app, nil
}

// initDatabase initializes the database connection
func (a *App) initDatabase() error {
	db, err := database.NewConnection(a.config.Database, a.logger)
	if err != nil {
		return err
	}

	a.db = db

	// Run migrations if auto-migrate is enabled
	if a.config.Database.AutoMigrate {
		migrator := database.NewMigrator(db, a.logger)
		if err := migrator.Up(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return nil
}

// initRepositories initializes all repositories
func (a *App) initRepositories() (*database.Repositories, error) {
	return database.NewRepositories(a.db, a.logger), nil
}

// initRouter initializes the HTTP router
func (a *App) initRouter() error {
	// Create gin engine
	a.router = gin.New()

	// Apply middlewares
	mwConfig := middleware.DefaultMiddlewareConfig()
	mwConfig.RateLimitReqs = a.config.Security.RateLimitRequests
	mwConfig.RateLimitWindow = a.config.Security.RateLimitWindow

	middleware.ApplyMiddlewares(a.router, a.config, a.logger, mwConfig)

	// Setup routes
	router.SetupRoutes(a.router, a.handlers, a.config)

	return nil
}

// Router returns the HTTP router
func (a *App) Router() http.Handler {
	return a.router
}

// Close closes all application resources
func (a *App) Close() error {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error("Failed to close database connection", "error", err)
			return err
		}
	}
	return nil
}

// Health returns the application health status
func (a *App) Health() map[string]interface{} {
	health := map[string]interface{}{
		"status": "ok",
		"app": map[string]interface{}{
			"name":        a.config.App.Name,
			"version":     a.config.App.Version,
			"environment": a.config.App.Environment,
		},
		"features": map[string]interface{}{
			"metrics":    a.config.Features.EnableMetrics,
			"tracing":    a.config.Features.EnableTracing,
			"profiling":  a.config.Features.EnableProfiling,
			"swagger":    a.config.Features.EnableSwagger,
			"healthz":    a.config.Features.EnableHealthz,
			"playground": a.config.Features.EnablePlayground,
		},
	}

	// Check database health
	if a.db != nil {
		if err := a.db.Ping(); err != nil {
			health["database"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
			health["status"] = "degraded"
		} else {
			health["database"] = map[string]interface{}{
				"status": "ok",
				"driver": a.config.Database.Driver,
			}
		}
	}

	return health
}

// Ready returns true if the application is ready to serve requests
func (a *App) Ready() bool {
	if a.db != nil {
		return a.db.Ping() == nil
	}
	return true
}

// Live returns true if the application is live
func (a *App) Live() bool {
	// Basic liveness check
	return true
}
