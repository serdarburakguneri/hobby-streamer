package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

type AssetAppService interface {
}

type EventHandlers struct {
	appService AssetAppService
	producer   *events.Producer
	logger     *logger.Logger
}

func NewEventHandlers(app AssetAppService, producer *events.Producer, l *logger.Logger) *EventHandlers {
	return &EventHandlers{appService: app, producer: producer, logger: l}
}

func (h *EventHandlers) HandleAnalyzeJobCompleted(ctx context.Context, ev *events.Event) error {
	var payload messages.JobCompletionPayload
	if err := unmarshalEventData(h.logger, ev, &payload); err != nil {
		return err
	}

	//TODO : Call UpdateVideoMetadata
	return nil
}

func (h *EventHandlers) HandleHLSJobCompleted(ctx context.Context, ev *events.Event) error {
	var payload map[string]interface{}
	if err := unmarshalEventData(h.logger, ev, &payload); err != nil {
		return err
	}
	return h.appService.AddVideo(ctx, payload["assetId"].(string), payload["videoId"].(string), payload)
}

func (h *EventHandlers) HandleRawVideoUploaded(ctx context.Context, ev *events.Event) error {
	return nil
}

func (h *EventHandlers) HandleTranscodeHlsJobCompleted(ctx context.Context, ev *events.Event) error {
	return h.HandleHLSJobCompleted(ctx, ev)
}

func (h *EventHandlers) HandleTranscodeDashJobCompleted(ctx context.Context, ev *events.Event) error {
	return nil
}
