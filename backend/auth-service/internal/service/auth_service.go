package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/models"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AuthService struct {
	keycloakURL  string
	realm        string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	keyFunc      jwt.Keyfunc
	logger       *logger.Logger
}

func NewAuthService(keycloakURL, realm, clientID, clientSecret string, keyFunc jwt.Keyfunc) *AuthService {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &AuthService{
		keycloakURL:  keycloakURL,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
		keyFunc:      keyFunc,
		logger:       logger.Get().WithService("auth-service"),
	}
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", req.ClientID)
	data.Set("username", req.Username)
	data.Set("password", req.Password)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.realm)
	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("failed to make token request", err), map[string]interface{}{
			"operation": "login",
			"username":  req.Username,
			"client_id": req.ClientID,
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("token request failed with status", fmt.Errorf("status: %d", resp.StatusCode)), map[string]interface{}{
			"operation":   "login",
			"status_code": resp.StatusCode,
			"username":    req.Username,
			"client_id":   req.ClientID,
		})
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("failed to decode token response", err), map[string]interface{}{
			"operation": "login",
			"username":  req.Username,
			"client_id": req.ClientID,
		})
	}

	return models.NewToken(tokenResponse.AccessToken, tokenResponse.TokenType, tokenResponse.ExpiresIn, tokenResponse.RefreshToken)
}

func (s *AuthService) ValidateToken(ctx context.Context, req *models.TokenValidationRequest) (*models.TokenValidationResult, error) {
	tokenString := strings.TrimPrefix(req.Token, "Bearer ")

	if tokenString == "" {
		return &models.TokenValidationResult{
			IsValid: false,
			Message: "Token is empty",
		}, nil
	}

	expiresAt, err := s.extractExpirationFromToken(tokenString)
	if err != nil {
		return &models.TokenValidationResult{
			IsValid: false,
			Message: "Invalid token payload",
		}, nil
	}

	if time.Now().After(expiresAt) {
		return &models.TokenValidationResult{
			IsValid:   false,
			Message:   "Token expired",
			ExpiresAt: expiresAt,
		}, nil
	}

	user, err := s.extractUserFromToken(tokenString)
	if err != nil {
		return &models.TokenValidationResult{
			IsValid:   false,
			Message:   "Invalid user data in token",
			ExpiresAt: expiresAt,
		}, nil
	}

	return &models.TokenValidationResult{
		IsValid:   true,
		User:      user,
		Message:   "",
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *models.TokenRefreshRequest) (*models.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("refresh_token", req.RefreshToken)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.realm)
	resp, err := s.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("failed to make refresh request", err), map[string]interface{}{
			"operation": "refresh_token",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("refresh request failed with status", fmt.Errorf("status: %d", resp.StatusCode)), map[string]interface{}{
			"operation":   "refresh_token",
			"status_code": resp.StatusCode,
		})
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("failed to decode refresh response", err), map[string]interface{}{
			"operation": "refresh_token",
		})
	}

	return models.NewToken(tokenResponse.AccessToken, tokenResponse.TokenType, tokenResponse.ExpiresIn, tokenResponse.RefreshToken)
}

func (s *AuthService) extractExpirationFromToken(tokenString string) (time.Time, error) {
	parsedToken, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		return time.Time{}, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "extract_expiration",
			"error":     err.Error(),
		})
	}

	if !parsedToken.Valid {
		return time.Time{}, pkgerrors.NewValidationError("invalid token", nil)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, pkgerrors.NewValidationError("invalid token claims", nil)
	}

	expClaim, exists := claims["exp"]
	if !exists {
		return time.Time{}, pkgerrors.NewValidationError("token missing expiration claim", nil)
	}

	var expTime time.Time
	switch exp := expClaim.(type) {
	case float64:
		expTime = time.Unix(int64(exp), 0)
	case int64:
		expTime = time.Unix(exp, 0)
	case string:
		expInt, err := time.Parse(time.RFC3339, exp)
		if err != nil {
			return time.Time{}, pkgerrors.WithContext(err, map[string]interface{}{
				"operation": "parse_expiration_string",
				"exp_value": exp,
			})
		}
		expTime = expInt
	default:
		return time.Time{}, pkgerrors.NewValidationError("invalid expiration claim format", nil)
	}

	return expTime, nil
}

func (s *AuthService) extractUserFromToken(tokenString string) (*models.User, error) {
	parsedToken, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "extract_user",
			"error":     err.Error(),
		})
	}

	if !parsedToken.Valid {
		return nil, pkgerrors.NewValidationError("invalid token", nil)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, pkgerrors.NewValidationError("invalid token claims", nil)
	}

	userIDStr := s.getStringClaim(claims, "sub")
	if userIDStr == "" {
		return nil, pkgerrors.NewValidationError("token missing user ID", nil)
	}

	usernameStr := s.getStringClaim(claims, "preferred_username")
	if usernameStr == "" {
		usernameStr = s.getStringClaim(claims, "username")
	}
	if usernameStr == "" {
		usernameStr = "unknown"
	}

	emailStr := s.getStringClaim(claims, "email")
	if emailStr == "" {
		emailStr = usernameStr + "@example.com"
	}

	var roleStrings []string
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if rolesInterface, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range rolesInterface {
				if roleStr, ok := role.(string); ok {
					roleStrings = append(roleStrings, roleStr)
				}
			}
		}
	}

	if len(roleStrings) == 0 {
		roleStrings = []string{"user"}
	}

	return models.NewUser(userIDStr, usernameStr, emailStr, roleStrings)
}

func (s *AuthService) getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
