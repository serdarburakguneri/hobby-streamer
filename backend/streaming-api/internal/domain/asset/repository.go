package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

type Repository interface {
	GetBySlug(ctx context.Context, slug valueobjects.Slug) (*entity.Asset, error)
	GetAll(ctx context.Context) ([]*entity.Asset, error)
	GetPublic(ctx context.Context) ([]*entity.Asset, error)
	GetByType(ctx context.Context, assetType valueobjects.AssetType) ([]*entity.Asset, error)
	GetByGenre(ctx context.Context, genre valueobjects.Genre) ([]*entity.Asset, error)
}
