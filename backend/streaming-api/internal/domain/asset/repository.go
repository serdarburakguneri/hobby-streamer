package asset

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id AssetID) (*Asset, error)
	GetBySlug(ctx context.Context, slug Slug) (*Asset, error)
	GetAll(ctx context.Context) ([]*Asset, error)
	GetPublic(ctx context.Context) ([]*Asset, error)
	GetByType(ctx context.Context, assetType AssetType) ([]*Asset, error)
	GetByGenre(ctx context.Context, genre Genre) ([]*Asset, error)
	GetByOwner(ctx context.Context, ownerID OwnerID) ([]*Asset, error)
	GetReady(ctx context.Context) ([]*Asset, error)
	GetPublished(ctx context.Context) ([]*Asset, error)
	Search(ctx context.Context, query string, filters *SearchFilters) ([]*Asset, error)
	GetRecommended(ctx context.Context, asset *Asset, limit int) ([]*Asset, error)
	GetStreamingInfo(ctx context.Context, asset *Asset, userID string, region string, userAge int) (*StreamingInfo, error)
}
