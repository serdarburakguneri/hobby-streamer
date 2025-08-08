package consumer

import (
	"context"

	cdn "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/cdn"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	domainentity "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
)

type AssetAppService interface {
	AddVideo(ctx context.Context, cmd commands.AddVideoCommand) error
	UpdateVideoMetadata(ctx context.Context, cmd commands.UpdateVideoMetadataCommand) error
	UpsertVideo(ctx context.Context, cmd commands.UpsertVideoCommand) (*domainentity.Asset, *domainentity.Video, error)
}

type EventHandlers struct {
	appService AssetAppService
	producer   *events.Producer
	cdn        cdn.Service
	logger     *logger.Logger
}

func NewEventHandlers(app AssetAppService, producer *events.Producer, cdnService cdn.Service, l *logger.Logger) *EventHandlers {
	return &EventHandlers{appService: app, producer: producer, cdn: cdnService, logger: l}
}
