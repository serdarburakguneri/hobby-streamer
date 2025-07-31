package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type ServiceClient struct {
	keycloakURL  string
	realm        string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	token        *ServiceToken
	tokenExpiry  time.Time
}

type ServiceToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

func NewServiceClient(keycloakURL, realm, clientID, clientSecret string) *ServiceClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &ServiceClient{
		keycloakURL:  keycloakURL,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
	}
}

func (s *ServiceClient) GetServiceToken(ctx context.Context) (string, error) {
	if s.token != nil && time.Now().Before(s.tokenExpiry) {
		return s.token.AccessToken, nil
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.realm)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", errors.NewInternalError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log := logger.WithService("service-client")
	log.Debug("Requesting service token", "token_url", tokenURL, "client_id", s.clientID)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", errors.NewInternalError("failed to request token", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.NewInternalError(fmt.Sprintf("token request failed with status: %d", resp.StatusCode), nil)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", errors.NewInternalError("failed to decode token response", err)
	}

	s.token = &ServiceToken{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
		Scope:       tokenResp.Scope,
	}

	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-30) * time.Second)

	log.Debug("Service token obtained successfully", "expires_in", tokenResp.ExpiresIn)

	return s.token.AccessToken, nil
}

func (s *ServiceClient) GetAuthorizationHeader(ctx context.Context) (string, error) {
	token, err := s.GetServiceToken(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", token), nil
}
