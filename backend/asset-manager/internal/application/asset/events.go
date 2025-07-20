package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
)

type EventPublisher interface {
	PublishAssetCreated(ctx context.Context, asset *asset.Asset) error
	PublishAssetUpdated(ctx context.Context, asset *asset.Asset) error
	PublishAssetDeleted(ctx context.Context, asset *asset.Asset) error
	PublishAssetPublished(ctx context.Context, asset *asset.Asset) error
	PublishVideoAdded(ctx context.Context, asset *asset.Asset, video *asset.Video) error
	PublishVideoRemoved(ctx context.Context, asset *asset.Asset, videoID string) error
	PublishVideoStatusUpdated(ctx context.Context, asset *asset.Asset, videoID string, status asset.VideoStatus) error
	PublishImageAdded(ctx context.Context, asset *asset.Asset, image asset.Image) error
	PublishImageRemoved(ctx context.Context, asset *asset.Asset, imageID string) error
}
