package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	graphql1 "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/interfaces/graphql"
)

type Resolver struct{}

// Assets is the resolver for the assets field.
func (r *bucketResolver) Assets(ctx context.Context, obj *graphql1.Bucket) ([]*graphql1.Asset, error) {
	panic("not implemented")
}

// CreateAsset is the resolver for the createAsset field.
func (r *mutationResolver) CreateAsset(ctx context.Context, input graphql1.CreateAssetInput) (*graphql1.Asset, error) {
	panic("not implemented")
}

// DeleteAsset is the resolver for the deleteAsset field.
func (r *mutationResolver) DeleteAsset(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

// UpdateAssetTitle is the resolver for the updateAssetTitle field.
func (r *mutationResolver) UpdateAssetTitle(ctx context.Context, id string, title string) (*graphql1.Asset, error) {
	panic("not implemented")
}

// UpdateAssetDescription is the resolver for the updateAssetDescription field.
func (r *mutationResolver) UpdateAssetDescription(ctx context.Context, id string, description string) (*graphql1.Asset, error) {
	panic("not implemented")
}

// SetAssetPublishRule is the resolver for the setAssetPublishRule field.
func (r *mutationResolver) SetAssetPublishRule(ctx context.Context, id string, rule graphql1.PublishRuleInput) (*graphql1.Asset, error) {
	panic("not implemented")
}

// ClearAssetPublishRule is the resolver for the clearAssetPublishRule field.
func (r *mutationResolver) ClearAssetPublishRule(ctx context.Context, id string) (*graphql1.Asset, error) {
	panic("not implemented")
}

// AddVideo is the resolver for the addVideo field.
func (r *mutationResolver) AddVideo(ctx context.Context, input graphql1.AddVideoInput) (*graphql1.Video, error) {
	panic("not implemented")
}

// DeleteVideo is the resolver for the deleteVideo field.
func (r *mutationResolver) DeleteVideo(ctx context.Context, assetID string, videoID string) (*graphql1.Asset, error) {
	panic("not implemented")
}

// RequestTranscode is the resolver for the requestTranscode field.
func (r *mutationResolver) RequestTranscode(ctx context.Context, assetID string, videoID string, format graphql1.VideoFormat) (bool, error) {
	panic("not implemented")
}

// CreateBucket is the resolver for the createBucket field.
func (r *mutationResolver) CreateBucket(ctx context.Context, input graphql1.BucketInput) (*graphql1.Bucket, error) {
	panic("not implemented")
}

// UpdateBucket is the resolver for the updateBucket field.
func (r *mutationResolver) UpdateBucket(ctx context.Context, id string, input graphql1.BucketInput) (*graphql1.Bucket, error) {
	panic("not implemented")
}

// DeleteBucket is the resolver for the deleteBucket field.
func (r *mutationResolver) DeleteBucket(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

// AddAssetToBucket is the resolver for the addAssetToBucket field.
func (r *mutationResolver) AddAssetToBucket(ctx context.Context, input graphql1.AddAssetToBucketInput) (bool, error) {
	panic("not implemented")
}

// RemoveAssetFromBucket is the resolver for the removeAssetFromBucket field.
func (r *mutationResolver) RemoveAssetFromBucket(ctx context.Context, input graphql1.RemoveAssetFromBucketInput) (bool, error) {
	panic("not implemented")
}

// AddImage is the resolver for the addImage field.
func (r *mutationResolver) AddImage(ctx context.Context, input graphql1.AddImageInput) (*graphql1.Asset, error) {
	panic("not implemented")
}

// Assets is the resolver for the assets field.
func (r *queryResolver) Assets(ctx context.Context, limit *int, offset *int) ([]*graphql1.Asset, error) {
	panic("not implemented")
}

// Asset is the resolver for the asset field.
func (r *queryResolver) Asset(ctx context.Context, id *string) (*graphql1.Asset, error) {
	panic("not implemented")
}

// ProcessingStatus is the resolver for the processingStatus field.
func (r *queryResolver) ProcessingStatus(ctx context.Context, assetID string, videoID string) (*graphql1.ProcessingStatus, error) {
	panic("not implemented")
}

// Buckets is the resolver for the buckets field.
func (r *queryResolver) Buckets(ctx context.Context, limit *int, nextKey *string) (*graphql1.BucketPage, error) {
	panic("not implemented")
}

// Bucket is the resolver for the bucket field.
func (r *queryResolver) Bucket(ctx context.Context, id *string) (*graphql1.Bucket, error) {
	panic("not implemented")
}

// BucketByKey is the resolver for the bucketByKey field.
func (r *queryResolver) BucketByKey(ctx context.Context, key string) (*graphql1.Bucket, error) {
	panic("not implemented")
}

// BucketsByOwner is the resolver for the bucketsByOwner field.
func (r *queryResolver) BucketsByOwner(ctx context.Context, ownerID string, limit *int, nextKey *string) (*graphql1.BucketPage, error) {
	panic("not implemented")
}

// SearchBuckets is the resolver for the searchBuckets field.
func (r *queryResolver) SearchBuckets(ctx context.Context, query string, limit *int, nextKey *string) (*graphql1.BucketPage, error) {
	panic("not implemented")
}

// SearchAssets is the resolver for the searchAssets field.
func (r *queryResolver) SearchAssets(ctx context.Context, query string, limit *int, offset *int) ([]*graphql1.Asset, error) {
	panic("not implemented")
}

// Bucket returns graphql1.BucketResolver implementation.
func (r *Resolver) Bucket() graphql1.BucketResolver { return &bucketResolver{r} }

// Mutation returns graphql1.MutationResolver implementation.
func (r *Resolver) Mutation() graphql1.MutationResolver { return &mutationResolver{r} }

// Query returns graphql1.QueryResolver implementation.
func (r *Resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type bucketResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct{}
*/
