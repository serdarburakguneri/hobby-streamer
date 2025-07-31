package models

import (
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

func NewLoginRequest(username, password, clientID string) (*LoginRequest, error) {
	if username == "" {
		return nil, pkgerrors.NewValidationError("username is required", nil)
	}
	if password == "" {
		return nil, pkgerrors.NewValidationError("password is required", nil)
	}
	if clientID == "" {
		return nil, pkgerrors.NewValidationError("client_id is required", nil)
	}

	return &LoginRequest{
		Username: username,
		Password: password,
		ClientID: clientID,
	}, nil
}

type TokenValidationRequest struct {
	Token string `json:"token"`
}

func NewTokenValidationRequest(token string) (*TokenValidationRequest, error) {
	if token == "" {
		return nil, pkgerrors.NewValidationError("token is required", nil)
	}

	return &TokenValidationRequest{
		Token: token,
	}, nil
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewTokenRefreshRequest(refreshToken string) (*TokenRefreshRequest, error) {
	if refreshToken == "" {
		return nil, pkgerrors.NewValidationError("refresh_token is required", nil)
	}

	return &TokenRefreshRequest{
		RefreshToken: refreshToken,
	}, nil
}
