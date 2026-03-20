# Go Base App

A comprehensive, production-ready Go application framework for building CLI tools, APIs, and microservices. This project follows clean architecture principles, modern development practices, and includes everything needed to build secure, scalable, and maintainable applications.

## Features

### Core Architecture
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Dependency Injection**: Modular and testable code structure
- **Configuration Management**: Environment-based configuration with validation
- **Structured Logging**: JSON logging with contextual information
- **Error Handling**: Comprehensive error types and handling

### Security & Compliance
- **12-Factor App Compliant**: Environment-based configuration
- **Security Headers**: CORS, XSS protection, content security policy
- **Rate Limiting**: Configurable request rate limiting
- **Input Validation**: Request validation with detailed error messages
- **Password Hashing**: Secure bcrypt password hashing

### Database & Persistence
- **Multi-Database Support**: SQLite, PostgreSQL, MySQL
- **Migration System**: Automated database migrations
- **Connection Pooling**: Optimized database connections
- **Health Checks**: Database connectivity monitoring

### HTTP API
- **RESTful API**: Well-structured REST endpoints
- **API Documentation**: Auto-generated Swagger/OpenAPI docs
- **Middleware Pipeline**: Logging, recovery, CORS, security
- **Request/Response**: JSON API with proper HTTP status codes

### Development & Operations
- **CLI Interface**: Comprehensive command-line tools
- **Hot Reload**: Development server with live reloading
- **Build System**: Makefile with common development tasks
- **Docker Support**: Containerization ready
- **Health Endpoints**: Kubernetes-compatible health checks

### Observability
- **Health Monitoring**: Application and database health checks
- **Metrics Ready**: Prometheus metrics endpoints
- **Graceful Shutdown**: Clean application termination

## Quick Start

### Prerequisites
- Go 1.21 or later
- Make (optional, for convenience commands)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ambiyansyah-risyal/go-base-app.git
cd go-base-app
```

2. Install dependencies:
```bash
go mod download
```

3. Build applications:
```bash
make build
# or manually:
# go build -o bin/go-base-app-api ./cmd/api
# go build -o bin/go-base-app-cli ./cmd/cli
```

### Configuration

Copy the example configuration:
```bash
cp .env.example .env
```

Edit `.env` file with your settings, or use environment variables:
```bash
export APP_NAME=my-app
export SERVER_PORT=8080
export DB_DRIVER=sqlite
export DB_DSN=./data/app.db
```

### Running the Application

#### API Server
```bash
# Using the binary
./bin/go-base-app-api

# Or using make
make run-api

# Or using go run
go run ./cmd/api
```

#### CLI Tools
```bash
# Show version
./bin/go-base-app-cli version

# Show current configuration
./bin/go-base-app-cli config show

# Validate configuration
./bin/go-base-app-cli config validate

# Run database migrations
./bin/go-base-app-cli migrate up
```

#### Development with Hot Reload
```bash
make dev
# This uses Air for hot reloading during development
```

## API Documentation

This application features **auto-generated OpenAPI/Swagger documentation** that's always up-to-date with the code. No more manually maintaining API docs!

### 📚 Interactive API Documentation

Once the server is running, you can explore the complete API documentation at:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Alternative endpoint**: http://localhost:8080/docs/index.html  
- **OpenAPI JSON**: http://localhost:8080/swagger/doc.json

![Swagger UI Screenshot](https://github.com/user-attachments/assets/f249136d-bf81-4eb6-89d0-aa4845ddf671)

### ✨ Features of Auto-Generated Documentation

- **Always Current**: Documentation is generated from code annotations, so it's never out of sync
- **Interactive Testing**: Test API endpoints directly from the browser
- **Complete Schemas**: Full request/response models with validation rules
- **Authentication Support**: Built-in support for Bearer token authentication
- **Search & Filter**: Easy navigation through endpoints by tags

### 🔄 Regenerating Documentation

The documentation is automatically included when building the application. To manually regenerate:

```bash
# Generate swagger documentation
make docs

# Or manually using swag CLI
swag init -g ./cmd/api/main.go -o ./docs

# Build with updated docs
make docs-build

# Start server with updated docs
make docs-serve
```

### 📖 Available Endpoints Overview

The API provides the following main endpoint groups:

- **🔐 Authentication** (`/api/v1/auth/*`) - User registration, login, token management
- **👥 Users** (`/api/v1/users/*`) - User management with full CRUD operations  
- **🩺 Health** (`/health*`) - Application and database health monitoring
- **📚 Documentation** (`/docs`, `/swagger`) - This interactive documentation

For detailed endpoint specifications, request/response schemas, and testing, visit the Swagger UI when the server is running.

## Quick API Testing

### 🚀 Interactive Testing (Recommended)

The easiest way to test the API is through the **interactive Swagger UI**:

1. Start the server: `make run-api` or `./bin/go-base-app-api`
2. Open your browser to: http://localhost:8080/swagger/index.html
3. Click "Authorize" to add Bearer token if needed
4. Expand any endpoint and click "Try it out" to test directly

### 🔧 Command Line Examples

For automation or scripting, here are some curl examples:

#### Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe", 
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

#### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

#### Check application health:
```bash
curl http://localhost:8080/health
```

💡 **Tip**: All request/response schemas, validation rules, and examples are available in the interactive Swagger documentation!

## Configuration

The application supports configuration via:
1. Environment variables (12-factor compliant)
2. Configuration files (YAML/JSON)
3. Command-line flags

### Key Configuration Options

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `APP_NAME` | `go-base-app` | Application name |
| `APP_ENVIRONMENT` | `development` | Environment (development/staging/production) |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `SERVER_PORT` | `8080` | Server port |
| `DB_DRIVER` | `sqlite` | Database driver (sqlite/postgres/mysql) |
| `DB_DSN` | `./data/app.db` | Database connection string |
| `LOG_LEVEL` | `info` | Log level (debug/info/warn/error) |
| `JWT_SECRET` | `change-in-production` | JWT signing secret |

## Development

### Project Structure
```
├── cmd/                    # Application entry points
│   ├── api/               # HTTP API server
│   ├── cli/               # CLI application  
│   └── migrate/           # Migration tool
├── internal/              # Private application code
│   ├── app/               # Application setup and wiring
│   ├── domain/            # Domain entities and interfaces
│   ├── infrastructure/    # External dependencies (DB, HTTP)
│   └── usecase/          # Business logic
├── pkg/                   # Public libraries
│   ├── config/           # Configuration management
│   ├── logger/           # Structured logging
│   ├── middleware/       # HTTP middleware
│   └── validator/        # Input validation
├── configs/              # Configuration files
├── migrations/           # Database migrations
├── docs/                # API documentation
└── test/                # Test files
```

### Available Make Commands

```bash
make help              # Show available commands
make build             # Build all applications
make test              # Run tests
make lint              # Run linter
make dev               # Start development server
make docs              # Generate API documentation
make docs-build        # Generate docs and build applications
make docs-serve        # Generate docs and serve API with Swagger UI
make clean             # Clean build artifacts
make docker-build      # Build Docker image
```

### Adding New Features

1. **Add Domain Entity**: Create new entities in `internal/domain/entity/`
2. **Add Repository Interface**: Define interfaces in `internal/domain/repository/`
3. **Implement Repository**: Add implementation in `internal/infrastructure/database/`
4. **Add Use Case**: Implement business logic in `internal/usecase/`
5. **Add HTTP Handler**: Create handlers in `internal/infrastructure/http/handler/`
6. **Update Routes**: Add routes in `internal/infrastructure/http/router/`
7. **📝 Document API**: Add Swagger annotations to new endpoints (see existing handlers for examples)

#### 📝 API Documentation Guidelines

When adding new API endpoints, always include Swagger annotations:

```go
// CreateWidget godoc
// @Summary Create a new widget
// @Description Create a new widget with the provided data
// @Tags widgets
// @Accept json
// @Produce json
// @Param request body CreateWidgetRequest true "Widget data"
// @Success 201 {object} Widget "Widget created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /api/v1/widgets [post]
func (h *WidgetHandler) Create(c *gin.Context) {
    // Implementation...
}
```

The documentation is automatically regenerated during build, keeping it always in sync with your code!

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run integration tests
make test-integration
```

## Deployment

### Docker
```bash
# Build Docker image
make docker-build

# Run with Docker
docker run -p 8080:8080 go-base-app:latest
```

### Environment Variables for Production
```bash
export APP_ENVIRONMENT=production
export JWT_SECRET=your-secure-secret-key
export DB_DRIVER=postgres
export DB_DSN=postgres://user:pass@host:5432/dbname
export LOG_LEVEL=warn
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run linter: `make lint`
6. Run tests: `make test`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Architecture Principles

This application follows these architectural principles:

- **Clean Architecture**: Clear separation between domain, use cases, and infrastructure
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Single Responsibility**: Each component has a single reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Convention over Configuration**: Sensible defaults with override capabilities
- **12-Factor App**: Environment-based configuration and stateless design
- **Security by Default**: Secure defaults for production deployment

## Roadmap

- [ ] JWT Authentication Implementation
- [ ] Role-based Access Control (RBAC)
- [ ] API Rate Limiting per User
- [ ] Prometheus Metrics Integration
- [ ] Distributed Tracing with OpenTelemetry
- [ ] GraphQL API Support
- [ ] Event Sourcing
- [ ] Microservice Communication Patterns
- [ ] Kubernetes Deployment Manifests
- [ ] Comprehensive Integration Tests