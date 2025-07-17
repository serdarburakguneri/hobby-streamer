package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	appconfig "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/config"
)

func main() {
	configManager, err := config.NewManager("streaming-api")
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
	log.Info("Starting streaming-api service", "environment", cfg.Environment)

	appConfig, err := appconfig.NewAppConfig(configManager, secretsManager, log)
	if err != nil {
		log.WithError(err).Error("Failed to initialize application configuration")
		os.Exit(1)
	}
	defer appConfig.Close()

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      appConfig.HTTP.Handler,
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
