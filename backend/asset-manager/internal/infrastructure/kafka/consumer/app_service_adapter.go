package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/queries"
	domainentity "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
)

type AssetAppServiceAdapter struct {
	commandService *asset.CommandService
	queryService   *asset.QueryService
}

func NewAssetAppServiceAdapter(commandService *asset.CommandService, queryService *asset.QueryService) *AssetAppServiceAdapter {
	return &AssetAppServiceAdapter{
		commandService: commandService,
		queryService:   queryService,
	}
}

func (a *AssetAppServiceAdapter) GetAsset(ctx context.Context, query queries.GetAssetQuery) (*domainentity.Asset, error) {
	return a.queryService.GetAsset(ctx, query)
}

func (a *AssetAppServiceAdapter) AddVideo(ctx context.Context, cmd commands.AddVideoCommand) error {
	return a.commandService.AddVideo(ctx, cmd)
}

func (a *AssetAppServiceAdapter) UpdateVideoStatus(ctx context.Context, cmd commands.UpdateVideoStatusCommand) error {
	return a.commandService.UpdateVideoStatus(ctx, cmd)
}

func (a *AssetAppServiceAdapter) UpdateVideoMetadata(ctx context.Context, cmd commands.UpdateVideoMetadataCommand) error {
	return a.commandService.UpdateVideoMetadata(ctx, cmd)
}

func (a *AssetAppServiceAdapter) UpsertVideo(ctx context.Context, cmd commands.UpsertVideoCommand) (*domainentity.Asset, *domainentity.Video, error) {
	return a.commandService.UpsertVideo(ctx, cmd)
}
