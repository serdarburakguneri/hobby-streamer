package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/auth"
	httphandler "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/http"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func main() {
	logLevel := getLogLevel()
	logFormat := getEnv("LOG_FORMAT", "text")
	logger.Init(logLevel, logFormat)
	log := logger.WithService("auth-service")

	log.Info("Starting auth-service")

	keycloakURL := getEnv("KEYCLOAK_URL", "https://localhost:8443")
	realm := getEnv("KEYCLOAK_REALM", "hobby")
	clientID := getEnv("KEYCLOAK_CLIENT_ID", "asset-manager")
	clientSecret := getEnv("KEYCLOAK_CLIENT_SECRET", "")

	log.Debug("Keycloak configuration", "url", keycloakURL, "realm", realm, "client_id", clientID)

	authService := auth.NewService(keycloakURL, realm, clientID, clientSecret)

	router := httphandler.NewRouter(authService)

	handler := logger.RequestLoggingMiddleware(log)(httphandler.CORS(router))

	port := getEnv("PORT", "8080")
	log.Info("Auth-service ready", "port", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.WithError(err).Error("Server failed to start")
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getLogLevel() slog.Level {
	level := getEnv("LOG_LEVEL", "info")
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
