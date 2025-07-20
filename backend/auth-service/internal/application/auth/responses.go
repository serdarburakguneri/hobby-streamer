package auth

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/user"
)

type TokenValidationResponse struct {
	valid   bool
	user    *user.User
	message string
	roles   []string
}

func NewTokenValidationResponse(valid bool, user *user.User, message string, roles []string) *TokenValidationResponse {
	return &TokenValidationResponse{
		valid:   valid,
		user:    user,
		message: message,
		roles:   roles,
	}
}

func (r *TokenValidationResponse) Valid() bool {
	return r.valid
}

func (r *TokenValidationResponse) User() *user.User {
	return r.user
}

func (r *TokenValidationResponse) Message() string {
	return r.message
}

func (r *TokenValidationResponse) Roles() []string {
	return r.roles
}
