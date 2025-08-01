package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AssetEventConsumer struct {
	appService *AssetAppServiceAdapter
	producer   *events.Producer
	consumer   *events.Consumer
	handlers   *EventHandlers
	logger     *logger.Logger
}

func NewAssetEventConsumer(appService *AssetAppServiceAdapter, producer *events.Producer) *AssetEventConsumer {
	l := logger.WithService("asset-event-consumer")
	return &AssetEventConsumer{
		appService: appService,
		producer:   producer,
		logger:     l,
		handlers:   NewEventHandlers(appService, producer, l),
	}
}

func (c *AssetEventConsumer) Start(ctx context.Context, bootstrapServers string) error {
	cfg := events.DefaultConsumerConfig()
	cfg.BootstrapServers = []string{bootstrapServers}
	cfg.GroupID = events.AssetManagerGroupID
	cfg.Topics = []string{
		events.RawVideoUploadedTopic,
		events.AnalyzeJobCompletedTopic,
		events.HLSJobCompletedTopic,
		events.DASHJobCompletedTopic,
	}

	cons, err := events.NewConsumer(ctx, cfg)
	if err != nil {
		return err
	}

	cons.Subscribe(events.RawVideoUploadedTopic, c.handlers.HandleRawVideoUploaded)
	cons.Subscribe(events.AnalyzeJobCompletedTopic, c.handlers.HandleAnalyzeJobCompleted)
	cons.Subscribe(events.HLSJobCompletedTopic, c.handlers.HandleTranscodeHlsJobCompleted)
	cons.Subscribe(events.DASHJobCompletedTopic, c.handlers.HandleTranscodeDashJobCompleted)

	c.consumer = cons
	go func() { _ = cons.Start(ctx) }()
	return nil
}

func (c *AssetEventConsumer) Stop() error {
	if c.consumer != nil {
		return c.consumer.Stop()
	}
	return nil
}
