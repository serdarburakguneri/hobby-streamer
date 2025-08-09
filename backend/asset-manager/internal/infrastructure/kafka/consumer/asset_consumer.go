package consumer

import (
	"context"

	cdn "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/cdn"
	apppipeline "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/pipeline"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AssetEventConsumer struct {
	appService *AssetAppServiceAdapter
	producer   *events.Producer
	consumer   *events.Consumer
	handlers   *EventHandlers
	logger     *logger.Logger
	cdnService cdn.Service
	pipeline   *apppipeline.Service
}

func NewAssetEventConsumer(appService *AssetAppServiceAdapter, publisher interface {
	Publish(context.Context, string, *events.Event) error
}, cdnService cdn.Service, pipelineSvc *apppipeline.Service) *AssetEventConsumer {
	l := logger.WithService("asset-event-consumer")
	return &AssetEventConsumer{
		appService: appService,
		producer:   nil,
		logger:     l,
		cdnService: cdnService,
		pipeline:   pipelineSvc,
		handlers:   NewEventHandlers(appService, publisher, cdnService, pipelineSvc, l),
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
