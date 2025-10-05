package main

import (
	"fmt"
	"os"

	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/spf13/cobra"
)

// Version information (set by build)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

// Global flags
var (
	configFile  string
	environment string
	verbose     bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "go-base-app",
	Short: "Go Base App CLI - A comprehensive Go application framework",
	Long: `Go Base App is a comprehensive base application for building
CLI tools, APIs, and microservices with Go. It provides a solid foundation
with clean architecture, security, observability, and modern development practices.`,
	Version: fmt.Sprintf("%s (built at %s)", Version, BuildTime),
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.go-base-app.yaml)")
	rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "development", "environment (development, staging, production)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	// Set environment variable if provided via flag
	if environment != "" {
		os.Setenv("APP_ENVIRONMENT", environment)
	}

	// Initialize basic logging for CLI
	loggerCfg := config.LoggerConfig{
		Level:      "info",
		Format:     "text",
		Output:     "stdout",
		Structured: false,
	}

	if verbose {
		loggerCfg.Level = "debug"
	}

	logger.Init(loggerCfg)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go Base App\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", "1.24+")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  "Start the HTTP API server with the specified configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		fmt.Printf("Starting API server on %s:%d...\n", cfg.Server.Host, cfg.Server.Port)
		fmt.Println("Use Ctrl+C to stop the server")
		
		// This would normally start the API server
		// For now, we'll just show the command would work
		return fmt.Errorf("serve command not fully implemented yet - use 'go run cmd/api/main.go' instead")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  "Run database migrations (up, down, status, etc.)",
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		fmt.Printf("Running migrations on %s database...\n", cfg.Database.Driver)
		// Migration logic would go here
		return fmt.Errorf("migrate up command not fully implemented yet")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Rolling back migrations...")
		// Migration rollback logic would go here
		return fmt.Errorf("migrate down command not fully implemented yet")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Migration status:")
		// Migration status logic would go here
		return fmt.Errorf("migrate status command not fully implemented yet")
	},
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  "Create, list, update, and delete users",
}

func init() {
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userDeleteCmd)
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating user...")
		// User creation logic would go here
		return fmt.Errorf("user create command not fully implemented yet")
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing users...")
		// User listing logic would go here
		return fmt.Errorf("user list command not fully implemented yet")
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete [user-id]",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		fmt.Printf("Deleting user: %s\n", userID)
		// User deletion logic would go here
		return fmt.Errorf("user delete command not fully implemented yet")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  "View and validate configuration",
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		fmt.Println("Current Configuration:")
		fmt.Printf("App Name: %s\n", cfg.App.Name)
		fmt.Printf("Version: %s\n", cfg.App.Version)
		fmt.Printf("Environment: %s\n", cfg.App.Environment)
		fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
		fmt.Printf("Database: %s\n", cfg.Database.Driver)
		fmt.Printf("Log Level: %s\n", cfg.Logger.Level)
		return nil
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}

		fmt.Println("✅ Configuration is valid")
		fmt.Printf("Environment: %s\n", cfg.App.Environment)
		return nil
	},
}