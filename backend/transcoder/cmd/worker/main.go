package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/app"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/job"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/queue"
)

func main() {
	logLevel := getLogLevel()
	logFormat := getEnv("LOG_FORMAT", "text")
	logger.Init(logLevel, logFormat)
	log := logger.WithService("transcoder-worker")

	log.Info("Starting transcoder worker")

	ctx := context.Background()
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Error("SQS_QUEUE_URL environment variable is not set")
		os.Exit(1)
	}

	log.Debug("Queue configuration", "queue_url", queueURL)

	var statusProducer *sqs.Producer
	statusQueueURL := os.Getenv("STATUS_QUEUE_URL")
	if statusQueueURL != "" {
		var err error
		statusProducer, err = sqs.NewProducer(ctx, statusQueueURL)
		if err != nil {
			log.WithError(err).Error("Failed to create status SQS producer, continuing without status updates")
		} else {
			log.Info("Status SQS producer initialized successfully", "status_queue_url", statusQueueURL)
		}
	}

	r := job.NewRegistry()
	if statusProducer != nil {
		r.Register("analyze", job.NewAnalyzeRunnerWithStatusProducer(statusProducer))
	} else {
		r.Register("analyze", job.NewAnalyzeRunner())
	}
	r.Register("transcode-hls", job.NewTranscodeHLSRunner())
	r.Register("transcode-dash", job.NewTranscodeDASHRunner())
	d := app.NewDispatcher(r)

	log.Info("Job registry initialized", "job_types", []string{"analyze", "transcode-hls", "transcode-dash"})

	consumer, err := queue.NewSQSConsumer(ctx, queueURL)
	if err != nil {
		log.WithError(err).Error("Failed to create queue consumer")
		os.Exit(1)
	}

	log.Info("Starting SQS consumer")
	consumer.Start(ctx, d.HandleMessage)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getLogLevel() slog.Level {
	level := getEnv("LOG_LEVEL", "info")
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
