package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
	appjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/application/job"
	domainjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/infrastructure/kafka"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/infrastructure/storage"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/infrastructure/transcoding"
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

	bootstrapServers := dynamicCfg.GetStringFromComponent("kafka", "bootstrap_servers")
	maxMessageBytes := dynamicCfg.GetIntFromComponent("kafka", "max_message_bytes")

	completionProducerConfig := &events.ProducerConfig{
		BootstrapServers: []string{bootstrapServers},
		Source:           "transcoder",
		MaxMessageBytes:  maxMessageBytes,
	}

	completionProducer, err := events.NewProducer(ctx, completionProducerConfig)
	if err != nil {
		log.WithError(err).Error("Failed to create completion producer")
		os.Exit(1)
	}

	kafkaEventPublisher := kafka.NewKafkaEventPublisher(completionProducer)
	s3Client, err := s3.NewClient(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to create S3 client")
		os.Exit(1)
	}
	storageAdapter := storage.NewStorage(s3Client)
	transcoderRegistry := transcoding.NewRegistry(storageAdapter)
	jobDomainService := domainjob.NewDomainService(storageAdapter, transcoderRegistry, kafkaEventPublisher)
	jobAppService := appjob.NewApplicationService(jobDomainService, dynamicCfg)

	transcoderEventConsumer := kafka.NewTranscoderEventConsumer(jobAppService, completionProducer)

	if err := transcoderEventConsumer.Start(ctx, bootstrapServers); err != nil {
		log.WithError(err).Error("Failed to start Kafka consumer")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down transcoder worker...")
	if err := transcoderEventConsumer.Stop(); err != nil {
		log.WithError(err).Error("Failed to stop Kafka consumer")
	}
	log.Info("Transcoder worker stopped")
}
