package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

// MockAuthService implements AuthService for testing
type MockAuthService struct {
	loginFunc         func(ctx context.Context, req *LoginRequest) (*Token, error)
	validateTokenFunc func(ctx context.Context, tokenString string) (*TokenValidationResponse, error)
	refreshTokenFunc  func(ctx context.Context, refreshToken string) (*Token, error)
}

func (m *MockAuthService) Login(ctx context.Context, req *LoginRequest) (*Token, error) {
	return m.loginFunc(ctx, req)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, tokenString string) (*TokenValidationResponse, error) {
	return m.validateTokenFunc(ctx, tokenString)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {
	return m.refreshTokenFunc(ctx, refreshToken)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{
		loginFunc: func(ctx context.Context, req *LoginRequest) (*Token, error) {
			return &Token{
				AccessToken:  "test-access-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				RefreshToken: "test-refresh-token",
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "testpass",
		ClientID: "test-client",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Token
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessToken != "test-access-token" {
		t.Errorf("Expected access token 'test-access-token', got '%s'", response.AccessToken)
	}
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{}
	handler := NewAuthHandler(mockService)

	testCases := []struct {
		name     string
		request  LoginRequest
		expected int
	}{
		{
			name:     "missing username",
			request:  LoginRequest{Password: "testpass", ClientID: "test-client"},
			expected: http.StatusBadRequest,
		},
		{
			name:     "missing password",
			request:  LoginRequest{Username: "testuser", ClientID: "test-client"},
			expected: http.StatusBadRequest,
		},
		{
			name:     "missing client_id",
			request:  LoginRequest{Username: "testuser", Password: "testpass"},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			if w.Code != tc.expected {
				t.Errorf("Expected status %d, got %d", tc.expected, w.Code)
			}
		})
	}
}

func TestAuthHandler_Login_ServiceError(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{
		loginFunc: func(ctx context.Context, req *LoginRequest) (*Token, error) {
			return nil, apperrors.NewUnauthorizedError("authentication failed", nil)
		},
	}

	handler := NewAuthHandler(mockService)
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
		ClientID: "test-client",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_ValidateToken_Success(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{
		validateTokenFunc: func(ctx context.Context, tokenString string) (*TokenValidationResponse, error) {
			return &TokenValidationResponse{
				Valid: true,
				User: &User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
					Roles:    []string{"user"},
				},
				Roles: []string{"user"},
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)
	validateReq := TokenValidationRequest{Token: "Bearer test-token"}
	body, _ := json.Marshal(validateReq)
	req := httptest.NewRequest("POST", "/validate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ValidateToken(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response TokenValidationResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Valid {
		t.Error("Expected token to be valid")
	}

	if response.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response.User.Username)
	}
}

func TestAuthHandler_ValidateToken_InvalidRequest(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{}
	handler := NewAuthHandler(mockService)
	validateReq := TokenValidationRequest{Token: ""}
	body, _ := json.Marshal(validateReq)
	req := httptest.NewRequest("POST", "/validate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ValidateToken(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{
		refreshTokenFunc: func(ctx context.Context, refreshToken string) (*Token, error) {
			return &Token{
				AccessToken:  "new-access-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				RefreshToken: "new-refresh-token",
			}, nil
		},
	}

	handler := NewAuthHandler(mockService)
	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{RefreshToken: "old-refresh-token"}
	body, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RefreshToken(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Token
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessToken != "new-access-token" {
		t.Errorf("Expected access token 'new-access-token', got '%s'", response.AccessToken)
	}
}

func TestAuthHandler_RefreshToken_InvalidRequest(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{}
	handler := NewAuthHandler(mockService)
	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{RefreshToken: ""}
	body, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RefreshToken(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Health(t *testing.T) {
	logger.Init(slog.LevelError, "text")

	mockService := &MockAuthService{}
	handler := NewAuthHandler(mockService)
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}
