package graphql

import (
	"context"
	"strconv"

	assetAppQueries "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/queries"
	bucketAppQueries "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket/queries"
	bucketvo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
)

// Assets resolves Bucket.assets.
func (r *bucketResolver) Assets(ctx context.Context, obj *Bucket) ([]*Asset, error) {
	bid, err := bucketvo.NewBucketID(obj.ID)
	if err != nil {
		return nil, err
	}
	ids, err := r.bucketQueryService.GetBucketAssets(ctx, bucketAppQueries.GetBucketAssetsQuery{BucketID: *bid, Limit: nil})
	if err != nil {
		return nil, err
	}
	out := make([]*Asset, len(ids))
	for i, aid := range ids {
		a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: aid})
		if err != nil {
			return nil, err
		}
		out[i] = domainAssetToGraphQL(a)
	}
	return out, nil
}

// CreateBucket resolves createBucket mutation.
func (r *mutationResolver) CreateBucket(ctx context.Context, input BucketInput) (*Bucket, error) {
	cmd, err := MapCreateBucketInput(input)
	if err != nil {
		return nil, err
	}
	b, err := r.bucketCommandService.CreateBucket(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return domainBucketToGraphQL(b), nil
}

// UpdateBucket resolves updateBucket mutation.
func (r *mutationResolver) UpdateBucket(ctx context.Context, id string, input BucketInput) (*Bucket, error) {
	cmd, err := MapUpdateBucketInput(id, input)
	if err != nil {
		return nil, err
	}
	if err := r.bucketCommandService.UpdateBucket(ctx, cmd); err != nil {
		return nil, err
	}
	b, err := r.bucketQueryService.GetBucket(ctx, bucketAppQueries.GetBucketQuery{ID: cmd.ID})
	if err != nil {
		return nil, err
	}
	return domainBucketToGraphQL(b), nil
}

// DeleteBucket resolves deleteBucket mutation.
func (r *mutationResolver) DeleteBucket(ctx context.Context, id string) (bool, error) {
	cmd, err := MapDeleteBucketInput(id)
	if err != nil {
		return false, err
	}
	if err := r.bucketCommandService.DeleteBucket(ctx, cmd); err != nil {
		return false, err
	}
	return true, nil
}

// AddAssetToBucket resolves addAssetToBucket mutation.
func (r *mutationResolver) AddAssetToBucket(ctx context.Context, input AddAssetToBucketInput) (bool, error) {
	cmd, err := MapAddAssetToBucketInput(input)
	if err != nil {
		return false, err
	}
	if err := r.bucketCommandService.AddAssetToBucket(ctx, cmd); err != nil {
		return false, err
	}
	return true, nil
}

// RemoveAssetFromBucket resolves removeAssetFromBucket mutation.
func (r *mutationResolver) RemoveAssetFromBucket(ctx context.Context, input RemoveAssetFromBucketInput) (bool, error) {
	cmd, err := MapRemoveAssetFromBucketInput(input)
	if err != nil {
		return false, err
	}
	if err := r.bucketCommandService.RemoveAssetFromBucket(ctx, cmd); err != nil {
		return false, err
	}
	return true, nil
}

// Buckets resolves Query.buckets.
func (r *queryResolver) Buckets(ctx context.Context, limit *int, nextKey *string) (*BucketPage, error) {
	var offPtr *int
	if nextKey != nil {
		off, err := strconv.Atoi(*nextKey)
		if err != nil {
			return nil, err
		}
		offPtr = &off
	}
	q := bucketAppQueries.ListBucketsQuery{Limit: limit, Offset: offPtr}
	page, err := r.bucketQueryService.ListBucketsPage(ctx, q)
	if err != nil {
		return nil, err
	}
	return domainBucketPageToGraphQL(page), nil
}

// Bucket resolves Query.bucket.
func (r *queryResolver) Bucket(ctx context.Context, id *string) (*Bucket, error) {
	if id == nil {
		return nil, nil
	}
	bid, err := bucketvo.NewBucketID(*id)
	if err != nil {
		return nil, err
	}
	b, err := r.bucketQueryService.GetBucket(ctx, bucketAppQueries.GetBucketQuery{ID: *bid})
	if err != nil {
		return nil, err
	}
	return domainBucketToGraphQL(b), nil
}

// BucketByKey resolves Query.bucketByKey.
func (r *queryResolver) BucketByKey(ctx context.Context, key string) (*Bucket, error) {
	bk, err := bucketvo.NewBucketKey(key)
	if err != nil {
		return nil, err
	}
	b, err := r.bucketQueryService.GetBucketByKey(ctx, bucketAppQueries.GetBucketByKeyQuery{Key: *bk})
	if err != nil {
		return nil, err
	}
	return domainBucketToGraphQL(b), nil
}

// BucketsByOwner resolves Query.bucketsByOwner.
func (r *queryResolver) BucketsByOwner(ctx context.Context, ownerID string, limit *int, nextKey *string) (*BucketPage, error) {
	oid, err := bucketvo.NewOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	var offPtr *int
	if nextKey != nil {
		off, err := strconv.Atoi(*nextKey)
		if err != nil {
			return nil, err
		}
		offPtr = &off
	}
	q := bucketAppQueries.GetBucketsByOwnerQuery{OwnerID: *oid, Limit: limit, Offset: offPtr}
	page, err := r.bucketQueryService.GetBucketsByOwnerPage(ctx, q)
	if err != nil {
		return nil, err
	}
	return domainBucketPageToGraphQL(page), nil
}

// SearchBuckets resolves Query.searchBuckets.
func (r *queryResolver) SearchBuckets(ctx context.Context, query string, limit *int, nextKey *string) (*BucketPage, error) {
	var offPtr *int
	if nextKey != nil {
		off, err := strconv.Atoi(*nextKey)
		if err != nil {
			return nil, err
		}
		offPtr = &off
	}
	q := bucketAppQueries.SearchBucketsQuery{Query: query, Limit: limit, Offset: offPtr}
	page, err := r.bucketQueryService.SearchBucketsPage(ctx, q)
	if err != nil {
		return nil, err
	}
	return domainBucketPageToGraphQL(page), nil
}
