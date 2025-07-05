package main

import (
	"log"
	"net/http"
	"os"

	"github.com/serdarburakguneri/hobby-streamer/services/auth-service/internal/auth"
	httphandler "github.com/serdarburakguneri/hobby-streamer/services/auth-service/internal/http"
)

func main() {
	keycloakURL := getEnv("KEYCLOAK_URL", "http://localhost:8080")
	realm := getEnv("KEYCLOAK_REALM", "hobby")
	clientID := getEnv("KEYCLOAK_CLIENT_ID", "asset-manager")
	clientSecret := getEnv("KEYCLOAK_CLIENT_SECRET", "")

	authService := auth.NewService(keycloakURL, realm, clientID, clientSecret)

	router := httphandler.NewRouter(authService)

	log.Println("[auth-service] Listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
