# Multi-stage Docker build for Go Base App

# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o go-base-app-api ./cmd/api
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o go-base-app-cli ./cmd/cli

# Final stage
FROM alpine:latest

# Install ca-certificates and sqlite3
RUN apk --no-cache add ca-certificates sqlite

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=appuser:appgroup /app/go-base-app-api .
COPY --from=builder --chown=appuser:appgroup /app/go-base-app-cli .

# Copy configuration files
COPY --from=builder --chown=appuser:appgroup /app/configs ./configs
COPY --from=builder --chown=appuser:appgroup /app/.env.example .env.example

# Create data directory for SQLite
RUN mkdir -p data && chown -R appuser:appgroup data

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ./go-base-app-cli config validate || exit 1

# Run the application
CMD ["./go-base-app-api"]