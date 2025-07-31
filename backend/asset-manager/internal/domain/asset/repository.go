package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
)

type Saver interface {
	Save(ctx context.Context, asset *entity.Asset) error
	Update(ctx context.Context, asset *entity.Asset) error
	Delete(ctx context.Context, id valueobjects.AssetID) error
}

type Finder interface {
	FindByID(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error)
	FindBySlug(ctx context.Context, slug valueobjects.Slug) (*entity.Asset, error)
}

type Querier interface {
	List(ctx context.Context, limit *int, offset *int) ([]*entity.Asset, error)
	Search(ctx context.Context, query string, limit *int, offset *int) ([]*entity.Asset, error)
	FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID, limit *int, offset *int) ([]*entity.Asset, error)
	FindByParentID(ctx context.Context, parentID valueobjects.AssetID, limit *int, offset *int) ([]*entity.Asset, error)
	FindByType(ctx context.Context, assetType valueobjects.AssetType, limit *int, offset *int) ([]*entity.Asset, error)
	FindByGenre(ctx context.Context, genre valueobjects.Genre, limit *int, offset *int) ([]*entity.Asset, error)
	FindByTag(ctx context.Context, tag valueobjects.Tag, limit *int, offset *int) ([]*entity.Asset, error)
}

type Repository interface {
	Saver
	Finder
	Querier
}
