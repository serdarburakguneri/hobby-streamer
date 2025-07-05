package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*Token, error)
	ValidateToken(ctx context.Context, tokenString string) (*TokenValidationResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
}

type Service struct {
	keycloakURL  string
	realm        string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	keyFunc      jwt.Keyfunc // for JWT signature verification
}

var _ AuthService = (*Service)(nil)

// NewService creates a new Service with a default keyFunc (no verification)
func NewService(keycloakURL, realm, clientID, clientSecret string) *Service {
	return &Service{
		keycloakURL:  keycloakURL,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			return nil, errors.New("no keyfunc provided")
		},
	}
}

// NewServiceWithKeyFunc allows injecting a custom keyFunc (for tests or custom verification)
func NewServiceWithKeyFunc(keycloakURL, realm, clientID, clientSecret string, keyFunc jwt.Keyfunc) *Service {
	return &Service{
		keycloakURL:  keycloakURL,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		keyFunc:      keyFunc,
	}
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*Token, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", req.ClientID)
	data.Set("username", req.Username)
	data.Set("password", req.Password)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.realm)
	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return &token, nil
}

func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*TokenValidationResponse, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return &TokenValidationResponse{Valid: false, Message: "Invalid token format"}, nil
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return &TokenValidationResponse{Valid: false, Message: "Invalid token signature"}, nil
		}
		if strings.Contains(err.Error(), "token is expired") {
			return &TokenValidationResponse{Valid: false, Message: "Token expired"}, nil
		}
		return &TokenValidationResponse{Valid: false, Message: err.Error()}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &TokenValidationResponse{Valid: false, Message: "Invalid token claims"}, nil
	}

	user := &User{
		ID:       getStringClaim(claims, "sub"),
		Username: getStringClaim(claims, "preferred_username"),
		Email:    getStringClaim(claims, "email"),
	}

	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range roles {
				if roleStr, ok := role.(string); ok {
					user.Roles = append(user.Roles, roleStr)
				}
			}
		}
	}

	return &TokenValidationResponse{
		Valid: true,
		User:  user,
		Roles: user.Roles,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("refresh_token", refreshToken)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.realm)
	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to make refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh request failed with status: %d", resp.StatusCode)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return &token, nil
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
