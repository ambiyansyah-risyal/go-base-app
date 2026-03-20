package handler

import (
	"net/http"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/ambiyansyah-risyal/go-base-app/internal/usecase"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	logger      *logger.Logger
	validator   *validator.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase, log *logger.Logger, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		logger:      log,
		validator:   validator,
	}
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User         entity.UserPublic `json:"user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int64             `json:"expires_in"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST_BODY",
		})
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		if validationErr, ok := err.(entity.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"code":   "VALIDATION_FAILED",
				"errors": validationErr,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"code":  "INVALID_REQUEST_DATA",
		})
		return
	}

	// Attempt login
	user, tokens, err := h.authUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		clientIP := c.ClientIP()
		h.logger.AuthLog("login", "", clientIP, false, err.Error())

		if err == entity.ErrInvalidCredentials || err == entity.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
				"code":  "INVALID_CREDENTIALS",
			})
			return
		}

		if err == entity.ErrUserInactive {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Account is inactive",
				"code":  "USER_INACTIVE",
			})
			return
		}

		h.logger.ErrorWithContext("Login failed", err, map[string]any{
			"email":     req.Email,
			"client_ip": clientIP,
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Login failed",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Log successful login
	clientIP := c.ClientIP()
	h.logger.AuthLog("login", user.ID.String(), clientIP, true, "")

	c.JSON(http.StatusOK, LoginResponse{
		User:         user.ToPublic(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// Register godoc
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.CreateUserRequest true "User registration data"
// @Success 201 {object} LoginResponse "Registration successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST_BODY",
		})
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		if validationErr, ok := err.(entity.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"code":   "VALIDATION_FAILED",
				"errors": validationErr,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"code":  "INVALID_REQUEST_DATA",
		})
		return
	}

	// Register user
	user, tokens, err := h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		if err == entity.ErrUserEmailExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
				"code":  "EMAIL_EXISTS",
			})
			return
		}

		if err == entity.ErrUserUsernameExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Username already exists",
				"code":  "USERNAME_EXISTS",
			})
			return
		}

		h.logger.ErrorWithContext("Registration failed", err, map[string]any{
			"email":    req.Email,
			"username": req.Username,
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Registration failed",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Log successful registration
	clientIP := c.ClientIP()
	h.logger.AuthLog("register", user.ID.String(), clientIP, true, "")

	c.JSON(http.StatusCreated, LoginResponse{
		User:         user.ToPublic(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} LoginResponse "Token refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST_BODY",
		})
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Refresh token is required",
			"code":  "MISSING_REFRESH_TOKEN",
		})
		return
	}

	// Refresh tokens
	user, tokens, err := h.authUsecase.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
			"code":  "INVALID_REFRESH_TOKEN",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		User:         user.ToPublic(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate tokens
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a full implementation, you would:
	// 1. Get the user ID from the JWT token
	// 2. Invalidate the tokens (add to blacklist or remove from database)
	// 3. Log the logout event

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserPublic "Current user information"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// In a full implementation, you would:
	// 1. Extract user ID from the JWT token in the Authorization header
	// 2. Fetch the user from the database
	// 3. Return the user information

	c.JSON(http.StatusOK, gin.H{
		"message": "Get current user endpoint - requires authentication middleware",
		"note":    "This endpoint needs JWT authentication middleware to be implemented",
	})
}
