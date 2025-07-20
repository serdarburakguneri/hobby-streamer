package token

import (
	"errors"
	"regexp"
	"time"
)

type AccessToken struct {
	value string
}

func NewAccessToken(value string) (*AccessToken, error) {
	if value == "" {
		return nil, ErrInvalidAccessToken
	}

	tokenRegex := regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)
	if !tokenRegex.MatchString(value) {
		return nil, ErrInvalidAccessToken
	}

	return &AccessToken{value: value}, nil
}

func (t AccessToken) Value() string {
	return t.value
}

func (t AccessToken) Equals(other AccessToken) bool {
	return t.value == other.value
}

type TokenType struct {
	value string
}

func NewTokenType(value string) (*TokenType, error) {
	if value == "" {
		return nil, ErrInvalidTokenType
	}

	validTypes := map[string]bool{
		"Bearer": true,
		"Basic":  true,
	}

	if !validTypes[value] {
		return nil, ErrInvalidTokenType
	}

	return &TokenType{value: value}, nil
}

func (t TokenType) Value() string {
	return t.value
}

func (t TokenType) Equals(other TokenType) bool {
	return t.value == other.value
}

type ExpiresIn struct {
	value int
}

func NewExpiresIn(value int) (*ExpiresIn, error) {
	if value <= 0 {
		return nil, ErrInvalidExpiresIn
	}

	if value > 86400 { // Max 24 hours
		return nil, ErrInvalidExpiresIn
	}

	return &ExpiresIn{value: value}, nil
}

func (e ExpiresIn) Value() int {
	return e.value
}

func (e ExpiresIn) Equals(other ExpiresIn) bool {
	return e.value == other.value
}

type RefreshToken struct {
	value string
}

func NewRefreshToken(value string) (*RefreshToken, error) {
	if value == "" {
		return nil, ErrInvalidRefreshToken
	}

	tokenRegex := regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)
	if !tokenRegex.MatchString(value) {
		return nil, ErrInvalidRefreshToken
	}

	return &RefreshToken{value: value}, nil
}

func (r RefreshToken) Value() string {
	return r.value
}

func (r RefreshToken) Equals(other RefreshToken) bool {
	return r.value == other.value
}

type ExpiresAt struct {
	value time.Time
}

func NewExpiresAt(value time.Time) *ExpiresAt {
	return &ExpiresAt{value: value}
}

func (e ExpiresAt) Value() time.Time {
	return e.value
}

func (e ExpiresAt) IsExpired() bool {
	return time.Now().After(e.value)
}

func (e ExpiresAt) IsExpiringSoon(threshold time.Duration) bool {
	return time.Now().Add(threshold).After(e.value)
}

func (e ExpiresAt) TimeUntilExpiry() time.Duration {
	return e.value.Sub(time.Now())
}

var (
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidTokenType    = errors.New("invalid token type")
	ErrInvalidExpiresIn    = errors.New("invalid expires in value")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
