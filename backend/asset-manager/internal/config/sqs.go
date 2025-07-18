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
	ConsumerRegistry *sqs.ConsumerRegistry
}

func NewSQSConfig(ctx context.Context, configManager *config.Manager, assetService *asset.Service, log *logger.Logger) (*SQSConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	analyzeCompletedQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_completed_queue_url")
	hlsCompletedQueueURL := dynamicCfg.GetStringFromComponent("sqs", "hls_completed_queue_url")
	dashCompletedQueueURL := dynamicCfg.GetStringFromComponent("sqs", "dash_completed_queue_url")

	consumerRegistry := sqs.NewConsumerRegistry()

	messageRouter := consumer.NewMessageRouter(assetService)

	consumerRegistry.Register(analyzeCompletedQueueURL, messageRouter.HandleAnalyzeMessage)
	consumerRegistry.Register(hlsCompletedQueueURL, messageRouter.HandleHLSMessage)
	consumerRegistry.Register(dashCompletedQueueURL, messageRouter.HandleDASHMessage)

	go func() {
		if err := consumerRegistry.Start(ctx); err != nil {
			log.WithError(err).Error("Failed to start consumer registry")
		}
	}()

	log.Info("Message router initialized",
		"analyze_completed_queue_url", analyzeCompletedQueueURL,
		"hls_completed_queue_url", hlsCompletedQueueURL,
		"dash_completed_queue_url", dashCompletedQueueURL,
		"supported_messages", []string{"analyze-completed", "transcode-hls-completed", "transcode-dash-completed"})

	return &SQSConfig{
		ConsumerRegistry: consumerRegistry,
	}, nil
}
