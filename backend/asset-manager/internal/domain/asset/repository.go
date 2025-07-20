package asset

import "context"

type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, id string) (*Asset, error)
	FindBySlug(ctx context.Context, slug string) (*Asset, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Asset, error)
	List(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	Search(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error)
	Delete(ctx context.Context, id string) error
	FindParent(ctx context.Context, childID string) (*Asset, error)
	FindChildren(ctx context.Context, parentID string) ([]*Asset, error)
	FindByTypeAndGenre(ctx context.Context, assetType, genre string) ([]*Asset, error)
	Update(ctx context.Context, asset *Asset) error
}
