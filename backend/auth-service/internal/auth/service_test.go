package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testJWTKey = []byte("test-secret-key")

func testKeyFunc(token *jwt.Token) (interface{}, error) {
	return testJWTKey, nil
}

func TestService_Login_Success(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/protocol/openid-connect/token") {
			t.Errorf("Expected token endpoint, got %s", r.URL.Path)
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		// Verify form data
		if r.FormValue("grant_type") != "password" {
			t.Errorf("Expected grant_type 'password', got '%s'", r.FormValue("grant_type"))
		}
		if r.FormValue("client_id") != "test-client" {
			t.Errorf("Expected client_id 'test-client', got '%s'", r.FormValue("client_id"))
		}
		if r.FormValue("username") != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", r.FormValue("username"))
		}
		if r.FormValue("password") != "testpass" {
			t.Errorf("Expected password 'testpass', got '%s'", r.FormValue("password"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Token{
			AccessToken:  "test-access-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "test-refresh-token",
		})
	}))
	defer server.Close()

	service := NewService(server.URL, "test-realm", "test-client", "test-secret")
	req := &LoginRequest{
		Username: "testuser",
		Password: "testpass",
		ClientID: "test-client",
	}

	// Act
	token, err := service.Login(context.Background(), req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token.AccessToken != "test-access-token" {
		t.Errorf("Expected access token 'test-access-token', got '%s'", token.AccessToken)
	}
	if token.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", token.TokenType)
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("Expected expires in 3600, got %d", token.ExpiresIn)
	}
	if token.RefreshToken != "test-refresh-token" {
		t.Errorf("Expected refresh token 'test-refresh-token', got '%s'", token.RefreshToken)
	}
}

func TestService_Login_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid_grant"}`))
	}))
	defer server.Close()

	service := NewService(server.URL, "test-realm", "test-client", "test-secret")
	req := &LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
		ClientID: "test-client",
	}

	// Act
	token, err := service.Login(context.Background(), req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if token != nil {
		t.Errorf("Expected nil token, got %v", token)
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("Expected error to contain status 401, got '%s'", err.Error())
	}
}

func TestService_ValidateToken_ValidToken(t *testing.T) {
	// Arrange
	service := NewServiceWithKeyFunc("http://localhost:8080", "test-realm", "test-client", "test-secret", testKeyFunc)

	// Create a valid JWT token
	claims := jwt.MapClaims{
		"sub":                "user123",
		"preferred_username": "testuser",
		"email":              "test@example.com",
		"exp":                time.Now().Add(time.Hour).Unix(),
		"iat":                time.Now().Unix(),
		"realm_access": map[string]interface{}{
			"roles": []interface{}{"user", "admin"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(testJWTKey)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Act
	response, err := service.ValidateToken(context.Background(), "Bearer "+tokenString)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Valid {
		t.Errorf("Expected token to be valid, got message: %s", response.Message)
	}

	if response.User.ID != "user123" {
		t.Errorf("Expected user ID 'user123', got '%s'", response.User.ID)
	}
	if response.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response.User.Username)
	}
	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", response.User.Email)
	}

	expectedRoles := []string{"user", "admin"}
	if len(response.Roles) != len(expectedRoles) {
		t.Errorf("Expected %d roles, got %d", len(expectedRoles), len(response.Roles))
	}
	for i, role := range expectedRoles {
		if response.Roles[i] != role {
			t.Errorf("Expected role '%s' at index %d, got '%s'", role, i, response.Roles[i])
		}
	}
}

func TestService_ValidateToken_ExpiredToken(t *testing.T) {
	// Arrange
	service := NewServiceWithKeyFunc("http://localhost:8080", "test-realm", "test-client", "test-secret", testKeyFunc)

	// Create an expired JWT token
	claims := jwt.MapClaims{
		"sub":                "user123",
		"preferred_username": "testuser",
		"email":              "test@example.com",
		"exp":                time.Now().Add(-time.Hour).Unix(), // Expired
		"iat":                time.Now().Add(-2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(testJWTKey)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Act
	response, err := service.ValidateToken(context.Background(), "Bearer "+tokenString)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Valid {
		t.Error("Expected token to be invalid (expired)")
	}

	if response.Message != "Token expired" {
		t.Errorf("Expected message 'Token expired', got '%s'", response.Message)
	}
}

func TestService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	service := NewService("http://localhost:8080", "test-realm", "test-client", "test-secret")

	// Act
	response, err := service.ValidateToken(context.Background(), "Bearer invalid-token")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Valid {
		t.Error("Expected token to be invalid")
	}

	if response.Message != "Invalid token format" {
		t.Errorf("Expected message 'Invalid token format', got '%s'", response.Message)
	}
}

func TestService_ValidateToken_NoBearerPrefix(t *testing.T) {
	// Arrange
	service := NewServiceWithKeyFunc("http://localhost:8080", "test-realm", "test-client", "test-secret", testKeyFunc)

	// Create a valid JWT token
	claims := jwt.MapClaims{
		"sub":                "user123",
		"preferred_username": "testuser",
		"email":              "test@example.com",
		"exp":                time.Now().Add(time.Hour).Unix(),
		"iat":                time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(testJWTKey)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Act
	response, err := service.ValidateToken(context.Background(), tokenString)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Valid {
		t.Errorf("Expected token to be valid, got message: %s", response.Message)
	}
	if response.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response.User.Username)
	}
}

func TestService_RefreshToken_Success(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		// Verify form data
		if r.FormValue("grant_type") != "refresh_token" {
			t.Errorf("Expected grant_type 'refresh_token', got '%s'", r.FormValue("grant_type"))
		}
		if r.FormValue("client_id") != "test-client" {
			t.Errorf("Expected client_id 'test-client', got '%s'", r.FormValue("client_id"))
		}
		if r.FormValue("client_secret") != "test-secret" {
			t.Errorf("Expected client_secret 'test-secret', got '%s'", r.FormValue("client_secret"))
		}
		if r.FormValue("refresh_token") != "old-refresh-token" {
			t.Errorf("Expected refresh_token 'old-refresh-token', got '%s'", r.FormValue("refresh_token"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Token{
			AccessToken:  "new-access-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "new-refresh-token",
		})
	}))
	defer server.Close()

	service := NewService(server.URL, "test-realm", "test-client", "test-secret")

	// Act
	token, err := service.RefreshToken(context.Background(), "old-refresh-token")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token.AccessToken != "new-access-token" {
		t.Errorf("Expected access token 'new-access-token', got '%s'", token.AccessToken)
	}
	if token.RefreshToken != "new-refresh-token" {
		t.Errorf("Expected refresh token 'new-refresh-token', got '%s'", token.RefreshToken)
	}
}

func TestService_RefreshToken_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid_grant"}`))
	}))
	defer server.Close()

	service := NewService(server.URL, "test-realm", "test-client", "test-secret")

	// Act
	token, err := service.RefreshToken(context.Background(), "invalid-refresh-token")

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if token != nil {
		t.Errorf("Expected nil token, got %v", token)
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected error to contain status 400, got '%s'", err.Error())
	}
}

func TestGetStringClaim(t *testing.T) {
	// Arrange
	claims := jwt.MapClaims{
		"string_key": "string_value",
		"int_key":    123,
		"nil_key":    nil,
	}

	// Act & Assert
	if result := getStringClaim(claims, "string_key"); result != "string_value" {
		t.Errorf("Expected 'string_value', got '%s'", result)
	}

	if result := getStringClaim(claims, "int_key"); result != "" {
		t.Errorf("Expected empty string for non-string value, got '%s'", result)
	}

	if result := getStringClaim(claims, "nil_key"); result != "" {
		t.Errorf("Expected empty string for nil value, got '%s'", result)
	}

	if result := getStringClaim(claims, "missing_key"); result != "" {
		t.Errorf("Expected empty string for missing key, got '%s'", result)
	}
}
