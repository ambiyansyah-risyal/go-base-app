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

## API Endpoints

### Health Checks
- `GET /health` - Application health status
- `GET /healthz` - Kubernetes health check
- `GET /livez` - Liveness probe
- `GET /readyz` - Readiness probe

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/me` - Get current user

### Users
- `GET /api/v1/users` - List users (with pagination/filtering)
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user (soft delete)

### Documentation
- `GET /docs` - Swagger UI documentation
- `GET /swagger` - Alternative Swagger endpoint

## Example Usage

### Register a new user:
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

### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### List users:
```bash
curl http://localhost:8080/api/v1/users?limit=10&offset=0
```

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
make clean             # Clean build artifacts
make docker-build      # Build Docker image
make docs              # Generate API documentation
```

### Adding New Features

1. **Add Domain Entity**: Create new entities in `internal/domain/entity/`
2. **Add Repository Interface**: Define interfaces in `internal/domain/repository/`
3. **Implement Repository**: Add implementation in `internal/infrastructure/database/`
4. **Add Use Case**: Implement business logic in `internal/usecase/`
5. **Add HTTP Handler**: Create handlers in `internal/infrastructure/http/handler/`
6. **Update Routes**: Add routes in `internal/infrastructure/http/router/`

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