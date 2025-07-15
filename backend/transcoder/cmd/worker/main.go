package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/job"
)

func main() {
	logFormat := getEnv("LOG_FORMAT", "text")
	logger.Init(slog.LevelInfo, logFormat)
	log := logger.WithService("transcoder-worker")

	log.Info("Starting transcoder worker")

	ctx := context.Background()
	queueURL := os.Getenv("TRANSCODER_QUEUE_URL")
	if queueURL == "" {
		log.Error("TRANSCODER_QUEUE_URL environment variable is not set")
		os.Exit(1)
	}

	log.Debug("Queue configuration", "queue_url", queueURL)

	analyzeQueueURL := os.Getenv("ANALYZE_QUEUE_URL")
	var analyzeProducer *sqs.Producer
	if analyzeQueueURL != "" {
		var err error
		analyzeProducer, err = sqs.NewProducer(ctx, analyzeQueueURL)
		if err != nil {
			log.WithError(err).Error("Failed to create analyze SQS producer, continuing without analyze completion messages")
		} else {
			log.Info("Analyze SQS producer initialized successfully", "analyze_queue_url", analyzeQueueURL)
		}
	}

	analyzeRunner := job.NewAnalyzeRunnerWithAnalyzeProducer(analyzeProducer)
	transcodeHLSRunner := job.NewTranscodeHLSRunnerWithAnalyzeProducer(analyzeProducer)
	transcodeDASHRunner := job.NewTranscodeDASHRunnerWithAnalyzeProducer(analyzeProducer)

	registry := sqs.NewConsumerRegistry()
	registry.Register(queueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
		log.Info("SQS message received", "msgType", msgType, "payload", payload)
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		switch msgType {
		case "analyze":
			return analyzeRunner.Run(ctx, payloadBytes)
		case "transcode-hls":
			return transcodeHLSRunner.Run(ctx, payloadBytes)
		case "transcode-dash":
			return transcodeDASHRunner.Run(ctx, payloadBytes)
		default:
			return nil
		}
	})

	log.Info("Job registry initialized", "job_types", []string{"analyze", "transcode-hls", "transcode-dash"})

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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
