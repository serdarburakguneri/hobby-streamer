package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.WithService("events-example")

	producerConfig := &events.ProducerConfig{
		BootstrapServers: []string{"localhost:9092"},
		Source:           "example-producer",
	}

	producer, err := events.NewProducer(ctx, producerConfig)
	if err != nil {
		log.WithError(err).Error("Failed to create producer")
		os.Exit(1)
	}
	defer producer.Close()

	consumerConfig := &events.ConsumerConfig{
		BootstrapServers: []string{"localhost:9092"},
		GroupID:          "example-group",
		Topics:           []string{events.AssetEventsTopic},
	}

	consumer, err := events.NewConsumer(ctx, consumerConfig)
	if err != nil {
		log.WithError(err).Error("Failed to create consumer")
		os.Exit(1)
	}
	defer consumer.Stop()

	consumer.Subscribe(events.AssetEventsTopic, func(ctx context.Context, event *events.Event) error {
		return nil
	})

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.WithError(err).Error("Consumer error")
		}
	}()

	time.Sleep(2 * time.Second)

	assetEvent := events.NewAssetCreatedEvent("asset-123", "test-asset", "Test Asset", "movie")
	assetEvent.SetCorrelationID("req-456")

	err = producer.SendEvent(ctx, events.AssetEventsTopic, assetEvent)
	if err != nil {
		log.WithError(err).Error("Failed to send event")
	}

	videoEvent := events.NewVideoAddedEvent("asset-123", "video-456", "Main Video", "hls")
	videoEvent.SetCorrelationID("req-456")
	videoEvent.SetCausationID(assetEvent.ID)

	err = producer.SendEvent(ctx, events.AssetEventsTopic, videoEvent)
	if err != nil {
		log.WithError(err).Error("Failed to send event")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down")
}
