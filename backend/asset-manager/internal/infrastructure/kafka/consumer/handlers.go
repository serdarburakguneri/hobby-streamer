package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
)

type AssetAppService interface {
	AddVideo(ctx context.Context, cmd commands.AddVideoCommand) error
	UpdateVideoMetadata(ctx context.Context, cmd commands.UpdateVideoMetadataCommand) error
}

type EventHandlers struct {
	appService AssetAppService
	producer   *events.Producer
	logger     *logger.Logger
}

func NewEventHandlers(app AssetAppService, producer *events.Producer, l *logger.Logger) *EventHandlers {
	return &EventHandlers{appService: app, producer: producer, logger: l}
}
