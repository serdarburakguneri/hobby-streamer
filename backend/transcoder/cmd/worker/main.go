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

	statusQueueURL := os.Getenv("STATUS_QUEUE_URL")
	var statusProducer *sqs.Producer
	if statusQueueURL != "" {
		var err error
		statusProducer, err = sqs.NewProducer(ctx, statusQueueURL)
		if err != nil {
			log.WithError(err).Error("Failed to create status SQS producer, continuing without status updates")
		} else {
			log.Info("Status SQS producer initialized successfully", "status_queue_url", statusQueueURL)
		}
	}

	analyzeRunner := job.NewAnalyzeRunnerWithStatusProducer(statusProducer)
	transcodeHLSRunner := job.NewTranscodeHLSRunner()
	transcodeDASHRunner := job.NewTranscodeDASHRunner()

	registry := sqs.NewConsumerRegistry()
	registry.Register(queueURL, func(ctx context.Context, msgType string, payload map[string]interface{}) error {
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

	// Wait for shutdown signal
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
