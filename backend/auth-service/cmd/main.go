package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	appauth "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/application/auth"
	infrahttp "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/infrastructure/http"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/infrastructure/keycloak"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
)

func main() {
	configManager, err := config.NewManager("auth-service")
	if err != nil {
		logger.Get().WithError(err).Error("Failed to initialize config")
		os.Exit(1)
	}
	defer configManager.Close()

	secretsManager := config.NewSecretsManager()
	secretsManager.LoadFromEnvironment()

	cfg := configManager.GetConfig()

	logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	log := logger.WithService(cfg.Service)
	log.Info("Starting auth-service", "environment", cfg.Environment)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return nil, errors.New("no keyfunc provided")
	}

	dynamicCfg := configManager.GetDynamicConfig()
	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	realm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	clientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")
	clientSecret := dynamicCfg.GetStringFromComponent("keycloak", "client_secret")

	keycloakRepo := keycloak.NewRepository(
		keycloakURL,
		realm,
		clientID,
		clientSecret,
		keyFunc,
	)

	authService := appauth.NewService(keycloakRepo)

	router := infrahttp.NewRouter(authService)
	handler := logger.RequestLoggingMiddleware(log)(router)

	securityMiddleware := func(next http.Handler) http.Handler {
		handler := next

		handler = security.SecurityHeadersMiddleware()(handler)
		handler = security.RateLimitMiddleware(cfg.Security.RateLimit.Requests, cfg.Security.RateLimit.Window)(handler)
		handler = security.CORSMiddleware(
			cfg.Security.CORS.AllowedOrigins,
			cfg.Security.CORS.AllowedMethods,
			cfg.Security.CORS.AllowedHeaders,
		)(handler)
		handler = security.InputValidationMiddleware()(handler)
		handler = security.LoggingMiddleware(log)(handler)

		return handler
	}

	handler = securityMiddleware(handler)
	handler = logger.CompressionMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("Starting HTTP server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
		os.Exit(1)
	}

	log.Info("Server exited")
}
