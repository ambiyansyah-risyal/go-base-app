package handler

import (
	"net/http"
	"strconv"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/ambiyansyah-risyal/go-base-app/internal/usecase"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/logger"
	"github.com/ambiyansyah-risyal/go-base-app/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUsecase usecase.UserUsecase
	logger      *logger.Logger
	validator   *validator.Validator
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase, log *logger.Logger, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		logger:      log,
		validator:   validator,
	}
}

// List godoc
// @Summary List users
// @Description Get list of users with pagination and filtering
// @Tags users
// @Accept json
// @Produce json
// @Param email query string false "Filter by email"
// @Param username query string false "Filter by username"
// @Param is_active query boolean false "Filter by active status"
// @Param limit query integer false "Limit number of results" default(20)
// @Param offset query integer false "Offset for pagination" default(0)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_desc query boolean false "Sort direction (true for desc)" default(true)
// @Success 200 {object} map[string]interface{} "List of users"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	// Parse query parameters
	filter := entity.DefaultUserFilter()

	if email := c.Query("email"); email != "" {
		filter.Email = email
	}

	if username := c.Query("username"); username != "" {
		filter.Username = username
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}

	if sortDescStr := c.Query("sort_desc"); sortDescStr != "" {
		if sortDesc, err := strconv.ParseBool(sortDescStr); err == nil {
			filter.SortDesc = sortDesc
		}
	}

	// Call use case
	users, total, err := h.userUsecase.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.ErrorWithContext("Failed to list users", err, map[string]any{
			"filter": filter,
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list users",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Convert to public format
	var publicUsers []entity.UserPublic
	for _, user := range users {
		publicUsers = append(publicUsers, user.ToPublic())
	}

	c.JSON(http.StatusOK, gin.H{
		"users": publicUsers,
		"meta": gin.H{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entity.UserPublic "User details"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	user, err := h.userUsecase.GetByID(c.Request.Context(), userID)
	if err != nil {
		if entity.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
			return
		}

		h.logger.ErrorWithContext("Failed to get user", err, map[string]any{
			"user_id": userID,
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToPublic())
}

// Update godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body entity.UpdateUserRequest true "Update user request"
// @Success 200 {object} entity.UserPublic "Updated user"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	var req entity.UpdateUserRequest
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

	user, err := h.userUsecase.Update(c.Request.Context(), userID, req)
	if err != nil {
		if entity.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
			return
		}

		h.logger.ErrorWithContext("Failed to update user", err, map[string]any{
			"user_id": userID,
			"request": req,
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToPublic())
}

// Delete godoc
// @Summary Delete user
// @Description Soft delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	err = h.userUsecase.Delete(c.Request.Context(), userID)
	if err != nil {
		if entity.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
			return
		}

		h.logger.ErrorWithContext("Failed to delete user", err, map[string]any{
			"user_id": userID,
		})

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
