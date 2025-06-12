package main

import (
	"context"
	"log"
	"os"

	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/app"
	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/job"
	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/queue"
)

func main() {
	ctx := context.Background()
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL is not set")
	}

	r := job.NewRegistry()
	r.Register("analyze", &job.AnalyzeRunner{})
	r.Register("transcode-hls", &job.TranscodeHLSRunner{})
	r.Register("transcode-dash", &job.TranscodeDASHRunner{})
	d := app.NewDispatcher(r)

	consumer, err := queue.NewSQSConsumer(ctx, queueURL)
	if err != nil {
		log.Fatalf("Failed to create queue consumer: %v", err)
	}

	consumer.Start(ctx, d.HandleMessage)
}