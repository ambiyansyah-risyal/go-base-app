package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Username  string     `json:"username" db:"username"`
	Password  string     `json:"-" db:"password"` // Never include in JSON responses
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewUser creates a new user instance
func NewUser(email, username, password, firstName, lastName string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsDeleted returns true if the user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// SoftDelete marks the user as deleted
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
	u.IsActive = false
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// UpdateInfo updates user information
func (u *User) UpdateInfo(firstName, lastName string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
}

// ToPublic returns a version of the user safe for public consumption
func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserPublic represents user data safe for public consumption
type UserPublic struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserFilter represents filters for user queries
type UserFilter struct {
	Email    string
	Username string
	IsActive *bool
	Limit    int
	Offset   int
	SortBy   string
	SortDesc bool
}

// DefaultUserFilter returns default filter values
func DefaultUserFilter() UserFilter {
	return UserFilter{
		Limit:    20,
		Offset:   0,
		SortBy:   "created_at",
		SortDesc: true,
	}
}