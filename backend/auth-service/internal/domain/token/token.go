package token

import (
	"errors"
	"time"
)

type Token struct {
	accessToken  AccessToken
	tokenType    TokenType
	expiresIn    ExpiresIn
	refreshToken RefreshToken
	expiresAt    ExpiresAt
	version      int
}

func NewToken(
	accessToken AccessToken,
	tokenType TokenType,
	expiresIn ExpiresIn,
	refreshToken RefreshToken,
	expiresAt ExpiresAt,
) *Token {
	return &Token{
		accessToken:  accessToken,
		tokenType:    tokenType,
		expiresIn:    expiresIn,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
		version:      1,
	}
}

func (t *Token) AccessToken() AccessToken {
	return t.accessToken
}

func (t *Token) TokenType() TokenType {
	return t.tokenType
}

func (t *Token) ExpiresIn() ExpiresIn {
	return t.expiresIn
}

func (t *Token) RefreshToken() RefreshToken {
	return t.refreshToken
}

func (t *Token) ExpiresAt() ExpiresAt {
	return t.expiresAt
}

func (t *Token) Version() int {
	return t.version
}

func (t *Token) IsExpired() bool {
	return t.expiresAt.IsExpired()
}

func (t *Token) IsExpiringSoon(threshold time.Duration) bool {
	return t.expiresAt.IsExpiringSoon(threshold)
}

func (t *Token) TimeUntilExpiry() time.Duration {
	return t.expiresAt.TimeUntilExpiry()
}

func (t *Token) IsValid() bool {
	return !t.IsExpired()
}

func (t *Token) HasRefreshToken() bool {
	return t.refreshToken.Value() != ""
}

func (t *Token) Refresh(newAccessToken AccessToken, newExpiresIn ExpiresIn, newExpiresAt ExpiresAt) error {
	if !t.HasRefreshToken() {
		return ErrNoRefreshToken
	}

	t.accessToken = newAccessToken
	t.expiresIn = newExpiresIn
	t.expiresAt = newExpiresAt
	t.version++
	return nil
}

func (t *Token) Revoke() {
	if t.IsValid() {
		t.expiresAt = *NewExpiresAt(time.Now().Add(-time.Second))
		t.version++
	}
}

var (
	ErrNoRefreshToken = errors.New("no refresh token available")
)
