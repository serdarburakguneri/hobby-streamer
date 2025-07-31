package neo4j

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/infrastructure/neo4j/bucket"
)

type BucketRepositoryAdapter struct {
	repo *bucket.Repository
}

func NewBucketRepositoryAdapter(repo *bucket.Repository) *BucketRepositoryAdapter {
	return &BucketRepositoryAdapter{repo: repo}
}

func (a *BucketRepositoryAdapter) Save(ctx context.Context, b *entity.Bucket) error {
	return a.repo.Create(ctx, b)
}

func (a *BucketRepositoryAdapter) FindByID(ctx context.Context, id valueobjects.BucketID) (*entity.Bucket, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *BucketRepositoryAdapter) FindByKey(ctx context.Context, key valueobjects.BucketKey) (*entity.Bucket, error) {
	return a.repo.GetByKey(ctx, key.Value())
}

func (a *BucketRepositoryAdapter) Exists(ctx context.Context, id valueobjects.BucketID) (bool, error) {
	return a.repo.Exists(ctx, id)
}

func (a *BucketRepositoryAdapter) ExistsByKey(ctx context.Context, key valueobjects.BucketKey) (bool, error) {
	return a.repo.ExistsByKey(ctx, key)
}

func (a *BucketRepositoryAdapter) Update(ctx context.Context, b *entity.Bucket) error {
	return a.repo.Update(ctx, b)
}

func (a *BucketRepositoryAdapter) Delete(ctx context.Context, id valueobjects.BucketID) error {
	return a.repo.Delete(ctx, id)
}

func (a *BucketRepositoryAdapter) List(ctx context.Context, limit *int, offset *int) ([]*entity.Bucket, error) {
	return a.repo.List(ctx, limit, offset)
}

func (a *BucketRepositoryAdapter) Search(ctx context.Context, query string, limit *int, offset *int) ([]*entity.Bucket, error) {
	return a.repo.Search(ctx, query, limit, offset)
}

func (a *BucketRepositoryAdapter) FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID, limit *int, offset *int) ([]*entity.Bucket, error) {
	page, err := a.repo.GetByOwnerID(ctx, ownerID.Value(), limit, map[string]interface{}{"offset": 0})
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (a *BucketRepositoryAdapter) FindByType(ctx context.Context, bucketType valueobjects.BucketType, limit *int, offset *int) ([]*entity.Bucket, error) {
	return a.repo.FindByType(ctx, bucketType, limit, offset)
}

func (a *BucketRepositoryAdapter) FindByStatus(ctx context.Context, status valueobjects.BucketStatus, limit *int, offset *int) ([]*entity.Bucket, error) {
	return a.repo.FindByStatus(ctx, status, limit, offset)
}

func (a *BucketRepositoryAdapter) AddAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) error {
	return a.repo.AddAsset(ctx, bucketID, assetID)
}

func (a *BucketRepositoryAdapter) RemoveAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) error {
	return a.repo.RemoveAsset(ctx, bucketID, assetID)
}

func (a *BucketRepositoryAdapter) GetAssetIDs(ctx context.Context, bucketID valueobjects.BucketID, limit *int, lastKey map[string]interface{}) ([]string, error) {
	return a.repo.GetAssetIDs(ctx, bucketID, limit, lastKey)
}

func (a *BucketRepositoryAdapter) HasAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) (bool, error) {
	return a.repo.HasAsset(ctx, bucketID, assetID)
}

func (a *BucketRepositoryAdapter) AssetCount(ctx context.Context, bucketID valueobjects.BucketID) (int, error) {
	return a.repo.AssetCount(ctx, bucketID)
}
