package config

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type SQSProducerConfig struct {
	AnalyzeCompletionProducer *sqs.Producer
	HLSCompletionProducer     *sqs.Producer
	DASHCompletionProducer    *sqs.Producer
}

func NewSQSProducerConfig(ctx context.Context, configManager *config.Manager, log *logger.Logger) (*SQSProducerConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	analyzeCompletionQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_completed_queue_url")
	hlsCompletionQueueURL := dynamicCfg.GetStringFromComponent("sqs", "hls_completed_queue_url")
	dashCompletionQueueURL := dynamicCfg.GetStringFromComponent("sqs", "dash_completed_queue_url")

	log.Info("SQS producer configuration",
		"analyze_completion_queue_url", analyzeCompletionQueueURL,
		"hls_completion_queue_url", hlsCompletionQueueURL,
		"dash_completion_queue_url", dashCompletionQueueURL)

	analyzeCompletionProducer, err := sqs.NewProducer(ctx, analyzeCompletionQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create analyze completion SQS producer")
		return nil, err
	}

	hlsCompletionProducer, err := sqs.NewProducer(ctx, hlsCompletionQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create HLS completion SQS producer")
		return nil, err
	}

	dashCompletionProducer, err := sqs.NewProducer(ctx, dashCompletionQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create DASH completion SQS producer")
		return nil, err
	}

	log.Info("SQS producers initialized successfully",
		"analyze_completion_queue_url", analyzeCompletionQueueURL,
		"hls_completion_queue_url", hlsCompletionQueueURL,
		"dash_completion_queue_url", dashCompletionQueueURL)

	return &SQSProducerConfig{
		AnalyzeCompletionProducer: analyzeCompletionProducer,
		HLSCompletionProducer:     hlsCompletionProducer,
		DASHCompletionProducer:    dashCompletionProducer,
	}, nil
}
