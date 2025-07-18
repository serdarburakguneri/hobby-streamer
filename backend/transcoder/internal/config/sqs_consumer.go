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

type SQSConsumerConfig struct {
	ConsumerRegistry *sqs.ConsumerRegistry
}

func NewSQSConsumerConfig(ctx context.Context, configManager *config.Manager, producerConfig *SQSProducerConfig, log *logger.Logger) (*SQSConsumerConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	hlsQueueURL := dynamicCfg.GetStringFromComponent("sqs", "hls_queue_url")
	dashQueueURL := dynamicCfg.GetStringFromComponent("sqs", "dash_queue_url")
	analyzeQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_queue_url")

	log.Info("SQS consumer configuration", "hls_queue_url", hlsQueueURL, "dash_queue_url", dashQueueURL, "analyze_queue_url", analyzeQueueURL)

	analyzeRunner := job.NewAnalyzeRunnerWithCompletionProducer(producerConfig.AnalyzeCompletionProducer)
	transcodeHLSRunner := job.NewTranscodeHLSRunnerWithCompletionProducer(producerConfig.HLSCompletionProducer)
	transcodeDASHRunner := job.NewTranscodeDASHRunnerWithCompletionProducer(producerConfig.DASHCompletionProducer)

	registry := sqs.NewConsumerRegistry()

	registry.Register(analyzeQueueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		log.Info("Analyze queue message received", "msgType", msgType, "payload", payload)
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if msgType != messages.MessageTypeAnalyze {
			log.Warn("Unknown message type in analyze queue", "msgType", msgType)
			return nil
		}
		return analyzeRunner.Run(ctx, payloadBytes)
	})

	registry.Register(hlsQueueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		log.Info("HLS queue message received", "msgType", msgType, "payload", payload)
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if msgType != messages.MessageTypeTranscodeHLS {
			log.Warn("Unknown message type in HLS queue", "msgType", msgType)
			return nil
		}
		return transcodeHLSRunner.Run(ctx, payloadBytes)
	})

	registry.Register(dashQueueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		log.Info("DASH queue message received", "msgType", msgType, "payload", payload)
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if msgType != messages.MessageTypeTranscodeDASH {
			log.Warn("Unknown message type in DASH queue", "msgType", msgType)
			return nil
		}
		return transcodeDASHRunner.Run(ctx, payloadBytes)
	})

	log.Info("SQS consumer registry initialized", "hls_queue_url", hlsQueueURL, "dash_queue_url", dashQueueURL, "analyze_queue_url", analyzeQueueURL, "job_types", []string{messages.MessageTypeAnalyze, messages.MessageTypeTranscodeHLS, messages.MessageTypeTranscodeDASH})

	return &SQSConsumerConfig{
		ConsumerRegistry: registry,
	}, nil
}
