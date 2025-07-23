package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func validJWT() string {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}

func TestTokenLifecycle(t *testing.T) {
	accessToken, _ := NewAccessToken(validJWT())
	tokenType, _ := NewTokenType("Bearer")
	expiresIn, _ := NewExpiresIn(3600)
	refreshToken, _ := NewRefreshToken(validJWT())
	expiresAt := NewExpiresAt(time.Now().Add(1 * time.Hour))

	t.Run("create and access fields", func(t *testing.T) {
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt)
		assert.Equal(t, *accessToken, tok.AccessToken())
		assert.Equal(t, *tokenType, tok.TokenType())
		assert.Equal(t, *expiresIn, tok.ExpiresIn())
		assert.Equal(t, *refreshToken, tok.RefreshToken())
		assert.Equal(t, *expiresAt, tok.ExpiresAt())
		assert.Equal(t, 1, tok.Version())
	})

	t.Run("is expired and is valid", func(t *testing.T) {
		expiredAt := NewExpiresAt(time.Now().Add(-1 * time.Hour))
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiredAt)
		assert.True(t, tok.IsExpired())
		assert.False(t, tok.IsValid())
	})

	t.Run("is expiring soon", func(t *testing.T) {
		expiresSoon := NewExpiresAt(time.Now().Add(30 * time.Second))
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresSoon)
		assert.True(t, tok.IsExpiringSoon(1*time.Minute))
		assert.False(t, tok.IsExpiringSoon(10*time.Second))
	})

	t.Run("time until expiry", func(t *testing.T) {
		expiresAt := NewExpiresAt(time.Now().Add(2 * time.Hour))
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt)
		delta := tok.TimeUntilExpiry()
		assert.Greater(t, delta.Seconds(), float64(7199))
	})

	t.Run("has refresh token", func(t *testing.T) {
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt)
		assert.True(t, tok.HasRefreshToken())
		noRefresh := NewToken(*accessToken, *tokenType, *expiresIn, RefreshToken{""}, *expiresAt)
		assert.False(t, noRefresh.HasRefreshToken())
	})

	t.Run("refresh with valid refresh token", func(t *testing.T) {
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt)
		newAccessToken, _ := NewAccessToken(validJWT())
		newExpiresIn, _ := NewExpiresIn(1800)
		newExpiresAt := NewExpiresAt(time.Now().Add(30 * time.Minute))
		oldVersion := tok.Version()
		err := tok.Refresh(*newAccessToken, *newExpiresIn, *newExpiresAt)
		assert.NoError(t, err)
		assert.Equal(t, *newAccessToken, tok.AccessToken())
		assert.Equal(t, *newExpiresIn, tok.ExpiresIn())
		assert.Equal(t, *newExpiresAt, tok.ExpiresAt())
		assert.Equal(t, oldVersion+1, tok.Version())
	})

	t.Run("refresh without refresh token returns error", func(t *testing.T) {
		tok := NewToken(*accessToken, *tokenType, *expiresIn, RefreshToken{""}, *expiresAt)
		newAccessToken, _ := NewAccessToken(validJWT())
		newExpiresIn, _ := NewExpiresIn(1800)
		newExpiresAt := NewExpiresAt(time.Now().Add(30 * time.Minute))
		err := tok.Refresh(*newAccessToken, *newExpiresIn, *newExpiresAt)
		assert.Error(t, err)
		assert.Equal(t, ErrNoRefreshToken, err)
	})

	t.Run("revoke token", func(t *testing.T) {
		tok := NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt)
		oldVersion := tok.Version()
		tok.Revoke()
		assert.True(t, tok.IsExpired())
		assert.Equal(t, oldVersion+1, tok.Version())
	})
}
