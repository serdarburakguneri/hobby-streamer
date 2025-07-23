package auth

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/token"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/user"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type LoginRequest struct {
	username user.Username
	password string
	clientID string
}

func NewLoginRequest(username, password, clientID string) (*LoginRequest, error) {
	usernameVO, err := user.NewUsername(username)
	if err != nil {
		return nil, err
	}
	if password == "" {
		return nil, ErrInvalidPassword
	}
	if clientID == "" {
		return nil, ErrInvalidClientID
	}
	return &LoginRequest{
		username: *usernameVO,
		password: password,
		clientID: clientID,
	}, nil
}

func (r *LoginRequest) Username() user.Username {
	return r.username
}

func (r *LoginRequest) Password() string {
	return r.password
}

func (r *LoginRequest) ClientID() string {
	return r.clientID
}

func (r *LoginRequest) Validate() error {
	if r.username.Value() == "" {
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
	token token.AccessToken
}

func NewTokenValidationRequest(tokenStr string) (*TokenValidationRequest, error) {
	tokenVO, err := token.NewAccessToken(tokenStr)
	if err != nil {
		return nil, err
	}
	return &TokenValidationRequest{
		token: *tokenVO,
	}, nil
}

func (r *TokenValidationRequest) Token() token.AccessToken {
	return r.token
}

func (r *TokenValidationRequest) Validate() error {
	if r.token.Value() == "" {
		return ErrInvalidToken
	}
	return nil
}

type TokenRefreshRequest struct {
	refreshToken token.RefreshToken
}

func NewTokenRefreshRequest(refreshTokenStr string) (*TokenRefreshRequest, error) {
	refreshTokenVO, err := token.NewRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, err
	}
	return &TokenRefreshRequest{
		refreshToken: *refreshTokenVO,
	}, nil
}

func (r *TokenRefreshRequest) RefreshToken() token.RefreshToken {
	return r.refreshToken
}

func (r *TokenRefreshRequest) Validate() error {
	if r.refreshToken.Value() == "" {
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
