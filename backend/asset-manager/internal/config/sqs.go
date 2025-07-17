package config

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/consumer"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type SQSConfig struct {
	Producer         *sqs.Producer
	ConsumerRegistry *sqs.ConsumerRegistry
}

func NewSQSConfig(ctx context.Context, configManager *config.Manager, assetService *asset.Service, log *logger.Logger) (*SQSConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	transcoderQueueURL := dynamicCfg.GetStringFromComponent("sqs", "transcoder_queue_url")
	sqsProducer, err := sqs.NewProducer(ctx, transcoderQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create SQS producer")
		return nil, err
	}

	analyzeQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_queue_url")
	consumerRegistry := sqs.NewConsumerRegistry()

	messageRouter := consumer.NewMessageRouter(assetService)
	consumerRegistry.Register(analyzeQueueURL, messageRouter.HandleMessage)

	go func() {
		if err := consumerRegistry.Start(ctx); err != nil {
			log.WithError(err).Error("Failed to start consumer registry")
		}
	}()

	log.Info("Message router initialized", "queue_url", analyzeQueueURL, "supported_messages", []string{"analyze-completed", "transcode-hls-completed", "transcode-dash-completed"})

	return &SQSConfig{
		Producer:         sqsProducer,
		ConsumerRegistry: consumerRegistry,
	}, nil
}
