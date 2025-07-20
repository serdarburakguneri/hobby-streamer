package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	appjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/application/job"
	domainjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
	infrasqs "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/infrastructure/sqs"
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
	dynamicCfg := configManager.GetDynamicConfig()

	logger.Init(logger.GetLogLevel(cfg.Log.Level), cfg.Log.Format)
	log := logger.WithService(cfg.Service)
	log.Info("Starting transcoder worker", "environment", cfg.Environment)

	ctx := context.Background()

	completionQueueURL := dynamicCfg.GetStringFromComponent("sqs", "completion_queue_url")

	completionProducer, err := sqs.NewProducer(ctx, completionQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create completion producer")
		os.Exit(1)
	}

	domainService := domainjob.NewJobDomainService()

	eventPublisher := infrasqs.NewEventPublisher(completionProducer)

	appService := appjob.NewApplicationService(domainService, eventPublisher)

	consumer := infrasqs.NewConsumer(appService)

	if err := consumer.RegisterQueues(configManager); err != nil {
		log.WithError(err).Error("Failed to register SQS queues")
		os.Exit(1)
	}

	log.Info("Starting SQS consumer")
	if err := consumer.Start(ctx); err != nil {
		log.WithError(err).Error("Failed to start consumer")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down transcoder worker...")
	consumer.Stop()
	log.Info("Transcoder worker stopped")
}
