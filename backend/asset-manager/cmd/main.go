package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appconfig "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func main() {
	configManager, err := config.NewManager("asset-manager")
	if err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}
	defer configManager.Close()

	secretsManager := config.NewSecretsManager()
	secretsManager.LoadFromEnvironment()

	cfg := configManager.GetConfig()
	logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	log := logger.WithService(cfg.Service)
	log.Info("Starting asset-manager service", "environment", cfg.Environment)

	ctx := context.Background()
	appConfig, err := appconfig.NewAppConfig(ctx, configManager, secretsManager, log)
	if err != nil {
		log.WithError(err).Error("Failed to initialize application configuration")
		os.Exit(1)
	}
	defer appConfig.Close()

	router := appConfig.GraphQL.Router
	router.Use(logger.CompressionMiddleware)
	router.Use(appConfig.Security.Middleware)
	router.Use(appConfig.Auth.Middleware)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("Starting server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
