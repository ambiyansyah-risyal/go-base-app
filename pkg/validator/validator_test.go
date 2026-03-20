package validator

import (
	"testing"

	"github.com/ambiyansyah-risyal/go-base-app/internal/domain/entity"
)

func TestNew(t *testing.T) {
	v := New()
	if v == nil {
		t.Error("New should not return nil")
	}
}

func TestValidateValidInput(t *testing.T) {
	v := New()

	user := entity.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	err := v.Validate(&user)
	if err != nil {
		t.Errorf("Validation should pass for valid input, got: %v", err)
	}
}

func TestValidateInvalidEmail(t *testing.T) {
	v := New()

	user := entity.CreateUserRequest{
		Email:     "invalid-email",
		Username:  "testuser",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	err := v.Validate(&user)
	if err == nil {
		t.Error("Validation should fail for invalid email")
	}
}

func TestValidateShortPassword(t *testing.T) {
	v := New()

	user := entity.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "short", // Too short (min=8)
		FirstName: "Test",
		LastName:  "User",
	}

	err := v.Validate(&user)
	if err == nil {
		t.Error("Validation should fail for password too short")
	}
}

func TestValidateShortUsername(t *testing.T) {
	v := New()

	user := entity.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "ab", // Too short (min=3)
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	err := v.Validate(&user)
	if err == nil {
		t.Error("Validation should fail for username too short")
	}
}

func TestValidateEmptyFields(t *testing.T) {
	v := New()

	user := entity.CreateUserRequest{
		Email:     "",
		Username:  "",
		Password:  "",
		FirstName: "",
		LastName:  "",
	}

	err := v.Validate(&user)
	if err == nil {
		t.Error("Validation should fail for empty required fields")
	}
}

func TestPasswordLength(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"SecurePass123!", true},
		{"12345678", true},      // Exactly 8 characters
		{"1234567", false},      // Too short (min=8)
		{"", false},             // Empty
	}

	for _, test := range tests {
		v := New()
		user := entity.CreateUserRequest{
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  test.password,
			FirstName: "Test",
			LastName:  "User",
		}

		err := v.Validate(&user)
		isValid := err == nil

		if isValid != test.valid {
			t.Errorf("Password %q: expected valid=%v, got valid=%v, error: %v", 
				test.password, test.valid, isValid, err)
		}
	}
}