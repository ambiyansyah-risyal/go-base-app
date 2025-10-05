# Go Base App Makefile
.PHONY: help build test lint clean install-tools dev migrate run-api run-cli deps tidy fmt vet

# Variables
APP_NAME := go-base-app
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

deps: ## Download dependencies
	go mod download
	go mod verify

tidy: ## Tidy up dependencies
	go mod tidy

# Build
build: ## Build all applications
	@echo "Building applications..."
	go build $(LDFLAGS) -o bin/$(APP_NAME)-api ./cmd/api
	go build $(LDFLAGS) -o bin/$(APP_NAME)-cli ./cmd/cli
	go build $(LDFLAGS) -o bin/$(APP_NAME)-migrate ./cmd/migrate

build-api: ## Build API server
	go build $(LDFLAGS) -o bin/$(APP_NAME)-api ./cmd/api

build-cli: ## Build CLI application
	go build $(LDFLAGS) -o bin/$(APP_NAME)-cli ./cmd/cli

build-migrate: ## Build migration tool
	go build $(LDFLAGS) -o bin/$(APP_NAME)-migrate ./cmd/migrate

# Development
dev: ## Start development server with hot reload
	air -c .air.toml

run-api: build-api ## Run API server
	./bin/$(APP_NAME)-api

run-cli: build-cli ## Run CLI application
	./bin/$(APP_NAME)-cli

migrate: build-migrate ## Run database migrations
	./bin/$(APP_NAME)-migrate up

# Testing
test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests
	go test -v -race -tags=integration ./test/...

coverage: test ## Generate test coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality
fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run linter
	golangci-lint run

# Documentation
docs: ## Generate API documentation
	swag init -g ./cmd/api/main.go -o ./docs

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker
docker-build: ## Build Docker image
	docker build -t $(APP_NAME):$(VERSION) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# Production
release: clean test lint build ## Build production release
	@echo "Release $(VERSION) built successfully"