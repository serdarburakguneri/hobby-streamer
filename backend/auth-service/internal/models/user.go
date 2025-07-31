package models

import (
	"regexp"
	"strings"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(id, username, email string, roles []string) (*User, error) {
	if err := validateUserID(id); err != nil {
		return nil, err
	}
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validateRoles(roles); err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Username:  strings.ToLower(username),
		Email:     strings.ToLower(email),
		Roles:     roles,
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func validateUserID(id string) error {
	if id == "" {
		return pkgerrors.NewValidationError("user ID cannot be empty", nil)
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" || len(username) < 3 || len(username) > 50 {
		return pkgerrors.NewValidationError("username must be between 3 and 50 characters", nil)
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return pkgerrors.NewValidationError("username contains invalid characters", nil)
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return pkgerrors.NewValidationError("email cannot be empty", nil)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return pkgerrors.NewValidationError("invalid email format", nil)
	}

	return nil
}

func validateRoles(roles []string) error {
	if len(roles) == 0 {
		return nil
	}

	roleRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	for _, role := range roles {
		if !roleRegex.MatchString(role) {
			return pkgerrors.NewValidationError("invalid role format", nil)
		}
	}

	return nil
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) HasAnyRole(roles []string) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}
