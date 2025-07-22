package auth

import (
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type LoginRequest struct {
	username string
	password string
	clientID string
}

func NewLoginRequest(username, password, clientID string) *LoginRequest {
	return &LoginRequest{
		username: username,
		password: password,
		clientID: clientID,
	}
}

func (r *LoginRequest) Username() string {
	return r.username
}

func (r *LoginRequest) Password() string {
	return r.password
}

func (r *LoginRequest) ClientID() string {
	return r.clientID
}

func (r *LoginRequest) Validate() error {
	if r.username == "" {
		return ErrInvalidUsername
	}
	if r.password == "" {
		return ErrInvalidPassword
	}
	if r.clientID == "" {
		return ErrInvalidClientID
	}
	return nil
}

type TokenValidationRequest struct {
	token string
}

func NewTokenValidationRequest(token string) *TokenValidationRequest {
	return &TokenValidationRequest{
		token: token,
	}
}

func (r *TokenValidationRequest) Token() string {
	return r.token
}

func (r *TokenValidationRequest) Validate() error {
	if r.token == "" {
		return ErrInvalidToken
	}
	return nil
}

type TokenRefreshRequest struct {
	refreshToken string
}

func NewTokenRefreshRequest(refreshToken string) *TokenRefreshRequest {
	return &TokenRefreshRequest{
		refreshToken: refreshToken,
	}
}

func (r *TokenRefreshRequest) RefreshToken() string {
	return r.refreshToken
}

func (r *TokenRefreshRequest) Validate() error {
	if r.refreshToken == "" {
		return ErrInvalidRefreshToken
	}
	return nil
}

var (
	ErrInvalidUsername     = pkgerrors.NewValidationError("invalid username", nil)
	ErrInvalidPassword     = pkgerrors.NewValidationError("invalid password", nil)
	ErrInvalidClientID     = pkgerrors.NewValidationError("invalid client ID", nil)
	ErrInvalidToken        = pkgerrors.NewValidationError("invalid token", nil)
	ErrInvalidRefreshToken = pkgerrors.NewValidationError("invalid refresh token", nil)
)
