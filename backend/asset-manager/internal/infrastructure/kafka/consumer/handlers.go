package consumer

import (
	"context"

	cdn "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/cdn"
	apppipeline "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/pipeline"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	domainentity "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
)

type AssetAppService interface {
	UpdateVideoMetadata(ctx context.Context, cmd commands.UpdateVideoMetadataCommand) error
	UpsertVideo(ctx context.Context, cmd commands.UpsertVideoCommand) (*domainentity.Asset, *domainentity.Video, error)
}

type Publisher interface {
	Publish(ctx context.Context, topic string, ev *events.Event) error
}

type EventHandlers struct {
	appService AssetAppService
	publisher  Publisher
	cdn        cdn.Service
	logger     *logger.Logger
	pipeline   *apppipeline.Service
}

func NewEventHandlers(app AssetAppService, publisher Publisher, cdnService cdn.Service, pipelineSvc *apppipeline.Service, l *logger.Logger) *EventHandlers {
	return &EventHandlers{appService: app, publisher: publisher, cdn: cdnService, pipeline: pipelineSvc, logger: l}
}
