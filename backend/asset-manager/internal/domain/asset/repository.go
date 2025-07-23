package asset

import "context"

type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, id AssetID) (*Asset, error)
	FindBySlug(ctx context.Context, slug Slug) (*Asset, error)
	FindByIDs(ctx context.Context, ids []AssetID) ([]*Asset, error)
	List(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	Search(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	Delete(ctx context.Context, id AssetID) error
	FindParent(ctx context.Context, childID AssetID) (*Asset, error)
	FindChildren(ctx context.Context, parentID AssetID) ([]*Asset, error)
	FindByTypeAndGenre(ctx context.Context, assetType *AssetType, genre *Genre) ([]*Asset, error)
	Update(ctx context.Context, asset *Asset) error
}
