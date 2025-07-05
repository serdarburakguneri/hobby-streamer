package auth

import "time"

// User represents a user in the system
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Enabled  bool     `json:"enabled"`
}

// Token represents a JWT token
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	Sub      string   `json:"sub"`
	Username string   `json:"preferred_username"`
	Email    string   `json:"email"`
	Roles    []string `json:"realm_access.roles"`
	Exp      int64    `json:"exp"`
	Iat      int64    `json:"iat"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

// TokenValidationRequest represents a token validation request
type TokenValidationRequest struct {
	Token string `json:"token"`
}

// TokenValidationResponse represents a token validation response
type TokenValidationResponse struct {
	Valid   bool     `json:"valid"`
	User    *User    `json:"user,omitempty"`
	Message string   `json:"message,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}
