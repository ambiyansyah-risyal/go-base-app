package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ambiyansyah-risyal/go-base-app/pkg/config"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Logger middleware for structured logging
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get status code
		statusCode := c.Writer.Status()

		// Format path with query parameters
		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		log.Request(
			c.Request.Method,
			path,
			c.Request.UserAgent(),
			clientIP,
			statusCode,
			latency.String(),
		)
	}
}

// Recovery middleware for panic recovery
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.ErrorWithContext("panic recovered", fmt.Errorf("%v", err), map[string]any{
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
					"ip":     c.ClientIP(),
				})

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"code":  "INTERNAL_ERROR",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS middleware
func CORS(cfg config.SecurityConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (basic)
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// Rate limiting middleware
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// Timeout middleware adds timeout to requests
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new context with timeout
		ctx := c.Request.Context()
		cancel := func() {}
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(c.Request.Context(), timeout)
		}
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ValidateContentType middleware validates Content-Type for POST/PUT requests
func ValidateContentType(allowedTypes ...string) gin.HandlerFunc {
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"application/json"}
	}

	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "POST" || method == "PUT" || method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type header is required",
					"code":  "MISSING_CONTENT_TYPE",
				})
				c.Abort()
				return
			}

			// Check if content type is allowed
			allowed := false
			for _, allowedType := range allowedTypes {
				if strings.Contains(strings.ToLower(contentType), strings.ToLower(allowedType)) {
					allowed = true
					break
				}
			}

			if !allowed {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": fmt.Sprintf("Unsupported Content-Type: %s", contentType),
					"code":  "UNSUPPORTED_CONTENT_TYPE",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// Health check middleware (skip logging for health endpoints)
func SkipHealthCheck() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/healthz", "/ready", "/live"},
	})
}

// Generate a simple request ID (in production, consider using UUIDs)
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// MiddlewareConfig holds middleware configuration
type MiddlewareConfig struct {
	EnableRateLimit   bool
	EnableTimeout     bool
	EnableRequestID   bool
	EnableContentType bool
	RateLimitReqs     int
	RateLimitWindow   time.Duration
	RequestTimeout    time.Duration
}

// ApplyMiddlewares applies all configured middlewares to the gin engine
func ApplyMiddlewares(r *gin.Engine, cfg *config.Config, log *logger.Logger, mwCfg MiddlewareConfig) {
	// Recovery middleware (should be first)
	r.Use(Recovery(log))

	// Security headers
	r.Use(SecurityHeaders())

	// CORS
	r.Use(CORS(cfg.Security))

	// Request ID
	if mwCfg.EnableRequestID {
		r.Use(RequestID())
	}

	// Rate limiting
	if mwCfg.EnableRateLimit {
		r.Use(RateLimit(mwCfg.RateLimitReqs, mwCfg.RateLimitWindow))
	}

	// Timeout
	if mwCfg.EnableTimeout && mwCfg.RequestTimeout > 0 {
		r.Use(Timeout(mwCfg.RequestTimeout))
	}

	// Content type validation
	if mwCfg.EnableContentType {
		r.Use(ValidateContentType())
	}

	// Logger (should be after other middlewares)
	r.Use(Logger(log))
}

// DefaultMiddlewareConfig returns default middleware configuration
func DefaultMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		EnableRateLimit:   true,
		EnableTimeout:     true,
		EnableRequestID:   true,
		EnableContentType: true,
		RateLimitReqs:     100,
		RateLimitWindow:   time.Minute,
		RequestTimeout:    30 * time.Second,
	}
}
