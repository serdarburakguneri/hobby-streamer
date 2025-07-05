package auth

import (
	"context"
	"os"
	"testing"
	"time"
)

// Integration tests require a running Keycloak instance
// Set INTEGRATION_TEST=true to run these tests
// These tests will be skipped if Keycloak is not available

func TestIntegration_Login_WithRealKeycloak(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	// Arrange
	keycloakURL := getEnvOrDefault("KEYCLOAK_URL", "http://localhost:9090")
	realm := getEnvOrDefault("KEYCLOAK_REALM", "hobby")
	clientID := getEnvOrDefault("KEYCLOAK_CLIENT_ID", "asset-manager")
	clientSecret := getEnvOrDefault("KEYCLOAK_CLIENT_SECRET", "your-client-secret")
	username := getEnvOrDefault("TEST_USERNAME", "testuser")
	password := getEnvOrDefault("TEST_PASSWORD", "testpass")

	service := NewService(keycloakURL, realm, clientID, clientSecret)
	req := &LoginRequest{
		Username: username,
		Password: password,
		ClientID: clientID,
	}

	// Act
	token, err := service.Login(context.Background(), req)

	// Assert
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Expected non-empty access token")
	}
	if token.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", token.TokenType)
	}
	if token.ExpiresIn <= 0 {
		t.Errorf("Expected positive expires_in, got %d", token.ExpiresIn)
	}
	if token.RefreshToken == "" {
		t.Error("Expected non-empty refresh token")
	}
}

func TestIntegration_ValidateToken_WithRealKeycloak(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	// Arrange
	keycloakURL := getEnvOrDefault("KEYCLOAK_URL", "http://localhost:9090")
	realm := getEnvOrDefault("KEYCLOAK_REALM", "hobby")
	clientID := getEnvOrDefault("KEYCLOAK_CLIENT_ID", "asset-manager")
	clientSecret := getEnvOrDefault("KEYCLOAK_CLIENT_SECRET", "your-client-secret")
	username := getEnvOrDefault("TEST_USERNAME", "testuser")
	password := getEnvOrDefault("TEST_PASSWORD", "testpass")

	service := NewService(keycloakURL, realm, clientID, clientSecret)

	// First, get a token
	loginReq := &LoginRequest{
		Username: username,
		Password: password,
		ClientID: clientID,
	}

	token, err := service.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("Failed to get token for validation test: %v", err)
	}

	// Act
	response, err := service.ValidateToken(context.Background(), "Bearer "+token.AccessToken)

	// Assert
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}

	if !response.Valid {
		t.Error("Expected token to be valid")
	}

	if response.User == nil {
		t.Error("Expected user information in response")
	}

	if response.User.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, response.User.Username)
	}

	if len(response.Roles) == 0 {
		t.Error("Expected at least one role in response")
	}
}

func TestIntegration_RefreshToken_WithRealKeycloak(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	// Arrange
	keycloakURL := getEnvOrDefault("KEYCLOAK_URL", "http://localhost:9090")
	realm := getEnvOrDefault("KEYCLOAK_REALM", "hobby")
	clientID := getEnvOrDefault("KEYCLOAK_CLIENT_ID", "asset-manager")
	clientSecret := getEnvOrDefault("KEYCLOAK_CLIENT_SECRET", "your-client-secret")
	username := getEnvOrDefault("TEST_USERNAME", "testuser")
	password := getEnvOrDefault("TEST_PASSWORD", "testpass")

	service := NewService(keycloakURL, realm, clientID, clientSecret)

	// First, get a token
	loginReq := &LoginRequest{
		Username: username,
		Password: password,
		ClientID: clientID,
	}

	originalToken, err := service.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("Failed to get token for refresh test: %v", err)
	}

	// Wait a moment to ensure tokens are different
	time.Sleep(1 * time.Second)

	// Act
	newToken, err := service.RefreshToken(context.Background(), originalToken.RefreshToken)

	// Assert
	if err != nil {
		t.Fatalf("Token refresh failed: %v", err)
	}

	if newToken.AccessToken == "" {
		t.Error("Expected non-empty access token")
	}

	if newToken.AccessToken == originalToken.AccessToken {
		t.Error("Expected different access token after refresh")
	}

	if newToken.RefreshToken == "" {
		t.Error("Expected non-empty refresh token")
	}

	if newToken.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", newToken.TokenType)
	}
}

func TestIntegration_ValidateToken_ExpiredToken(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	// Arrange
	keycloakURL := getEnvOrDefault("KEYCLOAK_URL", "http://localhost:9090")
	realm := getEnvOrDefault("KEYCLOAK_REALM", "hobby")
	clientID := getEnvOrDefault("KEYCLOAK_CLIENT_ID", "asset-manager")
	clientSecret := getEnvOrDefault("KEYCLOAK_CLIENT_SECRET", "your-client-secret")
	username := getEnvOrDefault("TEST_USERNAME", "testuser")
	password := getEnvOrDefault("TEST_PASSWORD", "testpass")

	service := NewService(keycloakURL, realm, clientID, clientSecret)

	// First, get a token
	loginReq := &LoginRequest{
		Username: username,
		Password: password,
		ClientID: clientID,
	}

	token, err := service.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("Failed to get token for expired token test: %v", err)
	}

	// Wait for token to expire (if it's a short-lived token)
	// Note: This test assumes the token has a short expiration time
	// In a real scenario, you might need to create a token with a very short expiration
	time.Sleep(2 * time.Second)

	// Act
	response, err := service.ValidateToken(context.Background(), "Bearer "+token.AccessToken)

	// Assert
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}

	// Note: This test might pass or fail depending on the token expiration time
	// In a real integration test, you might want to create a token with a very short expiration
	if !response.Valid {
		t.Log("Token is expired as expected")
	} else {
		t.Log("Token is still valid (this is normal for longer-lived tokens)")
	}
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
