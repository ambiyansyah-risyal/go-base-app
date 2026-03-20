// Package main provides the API server entry point
//
// @title           Go Base App API
// @version         1.0
// @description     A comprehensive Go base application for building APIs and microservices
// @termsOfService  https://github.com/ambiyansyah-risyal/go-base-app
//
// @contact.name   API Support
// @contact.url    https://github.com/ambiyansyah-risyal/go-base-app
// @contact.email  support@example.com
//
// @license.name  MIT
// @license.url   https://github.com/ambiyansyah-risyal/go-base-app/blob/main/LICENSE
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ambiyansyah-risyal/go-base-app/internal/app"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/gin-gonic/gin"

	// Import generated docs for swagger (conditional import)
	_ "github.com/ambiyansyah-risyal/go-base-app/docs"
)

// Version information (set by build)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.Logger)
	log := logger.Default()

	// Set gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create and initialize the application
	application, err := app.NewApp(cfg, log)
	if err != nil {
		log.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Setup HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      application.Router(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Log startup information
	features := map[string]bool{
		"metrics":    cfg.Features.EnableMetrics,
		"tracing":    cfg.Features.EnableTracing,
		"profiling":  cfg.Features.EnableProfiling,
		"swagger":    cfg.Features.EnableSwagger,
		"healthz":    cfg.Features.EnableHealthz,
		"playground": cfg.Features.EnablePlayground,
	}

	log.StartupLog(cfg.App.Name, Version, cfg.App.Environment, cfg.Server.Port, features)

	// Start server in a goroutine
	go func() {
		log.Info("Starting HTTP server",
			"host", cfg.Server.Host,
			"port", cfg.Server.Port,
			"version", Version,
			"build_time", BuildTime,
		)

		if cfg.Security.EnableHTTPS {
			if cfg.Security.TLSCertFile == "" || cfg.Security.TLSKeyFile == "" {
				log.Error("HTTPS enabled but TLS certificate or key file not specified")
				os.Exit(1)
			}

			log.Info("Starting HTTPS server",
				"cert_file", cfg.Security.TLSCertFile,
				"key_file", cfg.Security.TLSKeyFile,
			)

			if err := server.ListenAndServeTLS(cfg.Security.TLSCertFile, cfg.Security.TLSKeyFile); err != nil && err != http.ErrServerClosed {
				log.Error("Failed to start HTTPS server", "error", err)
				os.Exit(1)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error("Failed to start HTTP server", "error", err)
				os.Exit(1)
			}
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.ShutdownLog(cfg.App.Name, "received shutdown signal")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulStop)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	// Close application resources
	if err := application.Close(); err != nil {
		log.Error("Failed to close application resources", "error", err)
		os.Exit(1)
	}

	log.Info("Server gracefully stopped")
}
