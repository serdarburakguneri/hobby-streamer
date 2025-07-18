package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	appconfig "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/config"
)

func main() {
	configManager, err := config.NewManager("transcoder")
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
	log.Info("Starting transcoder worker", "environment", cfg.Environment)

	ctx := context.Background()
	appConfig, err := appconfig.NewAppConfig(ctx, configManager, log)
	if err != nil {
		log.WithError(err).Error("Failed to initialize application configuration")
		os.Exit(1)
	}

	log.Info("Starting SQS consumer")
	if err := appConfig.SQS.Consumer.ConsumerRegistry.Start(ctx); err != nil {
		log.WithError(err).Error("Failed to start consumer registry")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down transcoder worker...")
	appConfig.SQS.Consumer.ConsumerRegistry.Stop()
	log.Info("Transcoder worker stopped")
}
