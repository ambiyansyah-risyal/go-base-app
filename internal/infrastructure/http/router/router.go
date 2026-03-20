package router

import (
	"net/http"

	"github.com/ambiyansyah-risyal/go-base-app/internal/infrastructure/http/handler"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, h *handler.Handlers, cfg *config.Config) {
	// Health check routes (no versioning)
	if cfg.Features.EnableHealthz {
		setupHealthRoutes(r, h)
	}

	// API versioning
	v1 := r.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Auth.Login)
			auth.POST("/register", h.Auth.Register)
			auth.POST("/logout", h.Auth.Logout)
			auth.POST("/refresh", h.Auth.RefreshToken)
			auth.GET("/me", h.Auth.GetCurrentUser) // requires authentication
		}

		// User routes (protected)
		users := v1.Group("/users")
		// users.Use(middleware.Auth()) // Add authentication middleware
		{
			users.GET("", h.User.List)
			users.GET("/:id", h.User.GetByID)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", h.User.Delete)
		}
	}

	// Swagger documentation
	if cfg.Features.EnableSwagger {
		setupSwaggerRoutes(r)
	}

	// GraphQL playground (if enabled)
	if cfg.Features.EnablePlayground {
		setupPlaygroundRoutes(r)
	}

	// Metrics endpoint
	if cfg.Features.EnableMetrics {
		setupMetricsRoutes(r)
	}

	// Profiling endpoints (development only)
	if cfg.Features.EnableProfiling && cfg.IsDevelopment() {
		setupProfilingRoutes(r)
	}

	// Catch-all route for 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"code":    "NOT_FOUND",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})

	// Method not allowed handler
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":   "Method Not Allowed",
			"code":    "METHOD_NOT_ALLOWED",
			"message": "The request method is not supported for this resource",
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
		})
	})
}

// setupHealthRoutes sets up health check endpoints
func setupHealthRoutes(r *gin.Engine, h *handler.Handlers) {
	health := r.Group("/health")
	{
		health.GET("", h.Health.Health)
		health.GET("/live", h.Health.Live)
		health.GET("/ready", h.Health.Ready)
	}

	// Alternative health check endpoints
	r.GET("/healthz", h.Health.Health)
	r.GET("/livez", h.Health.Live)
	r.GET("/readyz", h.Health.Ready)
}

// setupSwaggerRoutes sets up Swagger documentation routes
func setupSwaggerRoutes(r *gin.Engine) {
	docs := r.Group("/docs")
	{
		docs.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Alternative swagger endpoints
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// setupPlaygroundRoutes sets up GraphQL playground (if using GraphQL)
func setupPlaygroundRoutes(r *gin.Engine) {
	r.GET("/playground", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "GraphQL playground would be available here",
			"note":    "GraphQL is not implemented in this base application",
		})
	})
}

// setupMetricsRoutes sets up Prometheus metrics endpoints
func setupMetricsRoutes(r *gin.Engine) {
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Prometheus metrics would be available here",
			"note":    "Metrics collection is not fully implemented yet",
		})
	})
}

// setupProfilingRoutes sets up pprof endpoints for profiling
func setupProfilingRoutes(r *gin.Engine) {
	pprof := r.Group("/debug/pprof")
	{
		pprof.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pprof profiling endpoints would be available here",
				"note":    "Profiling is not fully implemented yet",
			})
		})
	}
}
