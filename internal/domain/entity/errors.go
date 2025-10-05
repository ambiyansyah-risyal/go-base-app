package entity

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserEmailExists is returned when trying to create a user with an existing email
	ErrUserEmailExists = errors.New("user with email already exists")

	// ErrUserUsernameExists is returned when trying to create a user with an existing username
	ErrUserUsernameExists = errors.New("user with username already exists")

	// ErrInvalidCredentials is returned when authentication fails
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserInactive is returned when trying to authenticate an inactive user
	ErrUserInactive = errors.New("user account is inactive")

	// ErrUserDeleted is returned when trying to access a deleted user
	ErrUserDeleted = errors.New("user account has been deleted")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when user doesn't have permission
	ErrForbidden = errors.New("forbidden")

	// ErrInternalServer is returned for internal server errors
	ErrInternalServer = errors.New("internal server error")

	// ErrNotImplemented is returned for not implemented features
	ErrNotImplemented = errors.New("not implemented")
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "validation errors"
	}
	return fmt.Sprintf("validation errors: %d fields failed validation", len(e))
}

// Add adds a validation error
func (e *ValidationErrors) Add(field, message, value string) {
	*e = append(*e, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// BusinessError represents a domain business logic error
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e BusinessError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewBusinessError creates a new business error
func NewBusinessError(code, message string, details ...string) *BusinessError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &BusinessError{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

// Common business error codes
const (
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeUserEmailExists    = "USER_EMAIL_EXISTS"
	ErrCodeUserUsernameExists = "USER_USERNAME_EXISTS"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeUserInactive       = "USER_INACTIVE"
	ErrCodeUserDeleted        = "USER_DELETED"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeNotImplemented     = "NOT_IMPLEMENTED"
)

// IsNotFoundError checks if an error is a "not found" type error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr ValidationErrors
	return errors.As(err, &validationErr)
}

// IsBusinessError checks if an error is a business error
func IsBusinessError(err error) bool {
	var businessErr *BusinessError
	return errors.As(err, &businessErr)
}
