package keycloak

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
	appauth "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/application/auth"
	domauth "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/token"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Repository struct {
	keycloakURL            string
	realm                  string
	clientID               string
	clientSecret           string
	httpClient             *http.Client
	keyFunc                jwt.Keyfunc
	tokenValidationService domauth.TokenValidationService
	logger                 *logger.Logger
}

func NewRepository(keycloakURL, realm, clientID, clientSecret string, keyFunc jwt.Keyfunc) *Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec // Development only - disable TLS verification for local Keycloak
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	tokenValidationService := domauth.NewDomainTokenValidationService(keyFunc)

	return &Repository{
		keycloakURL:            keycloakURL,
		realm:                  realm,
		clientID:               clientID,
		clientSecret:           clientSecret,
		httpClient:             httpClient,
		keyFunc:                keyFunc,
		tokenValidationService: tokenValidationService,
		logger:                 logger.Get().WithService("keycloak-repository"),
	}
}

func (r *Repository) Login(ctx context.Context, req *appauth.LoginRequest) (*token.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", req.ClientID())
	data.Set("username", req.Username())
	data.Set("password", req.Password())

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", r.keycloakURL, r.realm)
	resp, err := r.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("failed to make token request", err), map[string]interface{}{
			"operation": "login",
			"username":  req.Username(),
			"client_id": req.ClientID(),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("token request failed with status", fmt.Errorf("status: %d", resp.StatusCode)), map[string]interface{}{
			"operation":   "login",
			"status_code": resp.StatusCode,
			"username":    req.Username(),
			"client_id":   req.ClientID(),
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
			"username":  req.Username(),
			"client_id": req.ClientID(),
		})
	}

	accessToken, err := token.NewAccessToken(tokenResponse.AccessToken)
	if err != nil {
		return nil, err
	}

	tokenType, err := token.NewTokenType(tokenResponse.TokenType)
	if err != nil {
		return nil, err
	}

	expiresIn, err := token.NewExpiresIn(tokenResponse.ExpiresIn)
	if err != nil {
		return nil, err
	}

	var refreshToken *token.RefreshToken
	if tokenResponse.RefreshToken != "" {
		refreshToken, err = token.NewRefreshToken(tokenResponse.RefreshToken)
		if err != nil {
			return nil, err
		}
	} else {
		refreshToken, _ = token.NewRefreshToken("")
	}

	expiresAt := token.NewExpiresAt(time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second))

	return token.NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt), nil
}

func (r *Repository) ValidateToken(ctx context.Context, req *appauth.TokenValidationRequest) (*appauth.TokenValidationResponse, error) {
	tokenString := strings.TrimPrefix(req.Token(), "Bearer ")

	validation, err := r.tokenValidationService.ValidateToken(tokenString)
	if err != nil {
		return nil, pkgerrors.WithContext(pkgerrors.NewInternalError("token validation failed", err), map[string]interface{}{
			"operation": "validate_token",
		})
	}

	if !validation.IsValid {
		return appauth.NewTokenValidationResponse(false, nil, validation.Message, nil), nil
	}

	// Convert domain roles to string slice for response
	var roleStrings []string
	for _, role := range validation.User.Roles().Values() {
		roleStrings = append(roleStrings, role.Value())
	}

	return appauth.NewTokenValidationResponse(true, validation.User, validation.Message, roleStrings), nil
}

func (r *Repository) RefreshToken(ctx context.Context, req *appauth.TokenRefreshRequest) (*token.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", r.clientID)
	data.Set("client_secret", r.clientSecret)
	data.Set("refresh_token", req.RefreshToken())

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", r.keycloakURL, r.realm)
	resp, err := r.httpClient.PostForm(tokenURL, data)
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

	accessToken, err := token.NewAccessToken(tokenResponse.AccessToken)
	if err != nil {
		return nil, err
	}

	tokenType, err := token.NewTokenType(tokenResponse.TokenType)
	if err != nil {
		return nil, err
	}

	expiresIn, err := token.NewExpiresIn(tokenResponse.ExpiresIn)
	if err != nil {
		return nil, err
	}

	var refreshToken *token.RefreshToken
	if tokenResponse.RefreshToken != "" {
		refreshToken, err = token.NewRefreshToken(tokenResponse.RefreshToken)
		if err != nil {
			return nil, err
		}
	} else {
		refreshToken, _ = token.NewRefreshToken("")
	}

	expiresAt := token.NewExpiresAt(time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second))

	return token.NewToken(*accessToken, *tokenType, *expiresIn, *refreshToken, *expiresAt), nil
}
