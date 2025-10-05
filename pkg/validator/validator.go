package validator

import (
	"reflect"
	"strings"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with custom functionality
type Validator struct {
	validator *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Use JSON tag names for validation errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	registerCustomValidators(v)

	return &Validator{
		validator: v,
	}
}

// Validate validates a struct and returns validation errors
func (v *Validator) Validate(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		var validationErrors entity.ValidationErrors

		if validatorErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validatorError := range validatorErrors {
				validationErrors.Add(
					validatorError.Field(),
					getErrorMessage(validatorError),
					validatorError.Value().(string),
				)
			}
		}

		if validationErrors.HasErrors() {
			return validationErrors
		}
	}
	return nil
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// RegisterCustomValidator registers a custom validation function
func (v *Validator) RegisterCustomValidator(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + fe.Param() + " characters long"
	case "max":
		return "Must be at most " + fe.Param() + " characters long"
	case "len":
		return "Must be exactly " + fe.Param() + " characters long"
	case "alpha":
		return "Must contain only alphabetic characters"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must contain only numeric characters"
	case "url":
		return "Must be a valid URL"
	case "uri":
		return "Must be a valid URI"
	case "uuid":
		return "Must be a valid UUID"
	case "uuid4":
		return "Must be a valid UUID v4"
	case "latitude":
		return "Must be a valid latitude coordinate"
	case "longitude":
		return "Must be a valid longitude coordinate"
	case "datetime":
		return "Must be a valid datetime"
	case "password":
		return "Must contain at least one uppercase letter, one lowercase letter, one digit, and one special character"
	case "username":
		return "Must contain only alphanumeric characters, underscores, and hyphens"
	case "phone":
		return "Must be a valid phone number"
	default:
		return "Invalid value"
	}
}

// registerCustomValidators registers custom validation functions
func registerCustomValidators(v *validator.Validate) {
	// Password validator: at least 8 chars, 1 upper, 1 lower, 1 digit, 1 special
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		hasUpper := false
		hasLower := false
		hasDigit := false
		hasSpecial := false

		for _, char := range password {
			switch {
			case char >= 'A' && char <= 'Z':
				hasUpper = true
			case char >= 'a' && char <= 'z':
				hasLower = true
			case char >= '0' && char <= '9':
				hasDigit = true
			default:
				hasSpecial = true
			}
		}

		return hasUpper && hasLower && hasDigit && hasSpecial
	})

	// Username validator: alphanumeric, underscore, hyphen only
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		for _, char := range username {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '_' || char == '-') {
				return false
			}
		}
		return true
	})

	// Phone validator: basic phone number validation
	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		// Simple regex-like validation for phone numbers
		// In production, use a proper phone number validation library
		if len(phone) < 10 || len(phone) > 15 {
			return false
		}

		for _, char := range phone {
			if !((char >= '0' && char <= '9') || char == '+' || char == '-' || char == ' ' || char == '(' || char == ')') {
				return false
			}
		}
		return true
	})
}

// ValidateCreateUserRequest validates create user request
func ValidateCreateUserRequest(req *entity.CreateUserRequest) error {
	v := New()
	return v.Validate(req)
}

// ValidateUpdateUserRequest validates update user request
func ValidateUpdateUserRequest(req *entity.UpdateUserRequest) error {
	v := New()
	return v.Validate(req)
}

// Common validation patterns
const (
	// EmailPattern is a basic email validation pattern
	EmailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// PasswordPattern is a strong password pattern
	PasswordPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`

	// UsernamePattern allows alphanumeric, underscore, and hyphen
	UsernamePattern = `^[a-zA-Z0-9_-]+$`

	// PhonePattern is a basic phone number pattern
	PhonePattern = `^[\+]?[1-9][\d]{0,15}$`
)

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	v := New()
	return v.ValidateVar(email, "email") == nil
}

// IsValidPassword checks if a password meets requirements
func IsValidPassword(password string) bool {
	v := New()
	return v.ValidateVar(password, "password") == nil
}

// IsValidUsername checks if a username is valid
func IsValidUsername(username string) bool {
	v := New()
	return v.ValidateVar(username, "username,min=3,max=50") == nil
}

// Default validator instance
var defaultValidator *Validator

// Init initializes the default validator
func Init() {
	defaultValidator = New()
}

// Default returns the default validator instance
func Default() *Validator {
	if defaultValidator == nil {
		Init()
	}
	return defaultValidator
}
