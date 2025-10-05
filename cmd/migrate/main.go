package main

import (
	"fmt"
	"os"

	"github.com/ambiyansyah-risyal/go-base-app/internal/infrastructure/database"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
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

	// Create database connection
	db, err := database.NewConnection(cfg.Database, log)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create migrator
	migrator := database.NewMigrator(db, log)

	// Determine command
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Error("Migration failed", "error", err)
			os.Exit(1)
		}
		log.Info("Migrations completed successfully")
	
	case "version":
		fmt.Printf("Go Base App Migration Tool\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
	
	default:
		fmt.Printf("Usage: %s [up|version]\n", os.Args[0])
		fmt.Printf("  up      Run pending migrations\n")
		fmt.Printf("  version Show version information\n")
		os.Exit(1)
	}
}