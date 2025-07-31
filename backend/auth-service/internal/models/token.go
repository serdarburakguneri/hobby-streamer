package models

import (
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

func NewToken(accessToken, tokenType string, expiresIn int, refreshToken string) (*Token, error) {
	if accessToken == "" {
		return nil, pkgerrors.NewValidationError("access token cannot be empty", nil)
	}
	if tokenType == "" {
		return nil, pkgerrors.NewValidationError("token type cannot be empty", nil)
	}
	if expiresIn <= 0 {
		return nil, pkgerrors.NewValidationError("expires in must be positive", nil)
	}

	return &Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
	}, nil
}

type TokenValidationResult struct {
	IsValid   bool      `json:"is_valid"`
	User      *User     `json:"user,omitempty"`
	Message   string    `json:"message,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
