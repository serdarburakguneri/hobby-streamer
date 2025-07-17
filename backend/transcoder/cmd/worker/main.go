package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/job"
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
	transcoderQueueURL := dynamicCfg.GetStringFromComponent("sqs", "transcoder_queue_url")
	analyzeQueueURL := dynamicCfg.GetStringFromComponent("sqs", "analyze_queue_url")

	log.Info("Queue configuration", "transcoder_queue_url", transcoderQueueURL, "analyze_queue_url", analyzeQueueURL)

	analyzeProducer, err := sqs.NewProducer(ctx, analyzeQueueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create analyze SQS producer")
		os.Exit(1)
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

	log.Info("Starting SQS consumer")
	if err := registry.Start(ctx); err != nil {
		log.WithError(err).Error("Failed to start consumer registry")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down transcoder worker...")
	registry.Stop()
	log.Info("Transcoder worker stopped")
}
