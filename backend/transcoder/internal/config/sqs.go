package config

import (
	"context"
	"encoding/json"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/job"
)

type SQSConfig struct {
	AnalyzeProducer  *sqs.Producer
	ConsumerRegistry *sqs.ConsumerRegistry
}

func NewSQSConfig(ctx context.Context, configManager *config.Manager, log *logger.Logger) (*SQSConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	transcoderQueueURL := dynamicCfg.GetStringFromComponent("sqs", "transcoder_queue_url")
	analyzeQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_queue_url")

	log.Info("Queue configuration", "transcoder_queue_url", transcoderQueueURL, "analyze_queue_url", analyzeQueueURL)

	analyzeProducer, err := sqs.NewProducer(ctx, analyzeQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create analyze SQS producer")
		return nil, err
	}
	log.Info("Analyze SQS producer initialized successfully", "analyze_queue_url", analyzeQueueURL)

	analyzeRunner := job.NewAnalyzeRunnerWithAnalyzeProducer(analyzeProducer)
	transcodeHLSRunner := job.NewTranscodeHLSRunnerWithAnalyzeProducer(analyzeProducer)
	transcodeDASHRunner := job.NewTranscodeDASHRunnerWithAnalyzeProducer(analyzeProducer)

	registry := sqs.NewConsumerRegistry()
	registry.Register(transcoderQueueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		log.Info("SQS message received", "msgType", msgType, "payload", payload)
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		switch msgType {
		case messages.MessageTypeAnalyze:
			return analyzeRunner.Run(ctx, payloadBytes)
		case messages.MessageTypeTranscodeHLS:
			return transcodeHLSRunner.Run(ctx, payloadBytes)
		case messages.MessageTypeTranscodeDASH:
			return transcodeDASHRunner.Run(ctx, payloadBytes)
		default:
			return nil
		}
	})

	log.Info("Job registry initialized", "job_types", []string{messages.MessageTypeAnalyze, messages.MessageTypeTranscodeHLS, messages.MessageTypeTranscodeDASH})

	return &SQSConfig{
		AnalyzeProducer:  analyzeProducer,
		ConsumerRegistry: registry,
	}, nil
}
