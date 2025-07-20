package auth

import "context"

type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (*User, error)
	HasRole(user *User, role string) bool
	HasAnyRole(user *User, roles []string) bool
	HasAllRoles(user *User, roles []string) bool
}

type ValidationResult struct {
	Valid bool
	User  *User
	Error string
}

type ServiceClientInterface interface {
	GetServiceToken(ctx context.Context) (string, error)
	GetAuthorizationHeader(ctx context.Context) (string, error)
}
