package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"
	"fmt"
	"runtime/debug"

	appasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	appbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Assets is the resolver for the assets field.
func (r *bucketResolver) Assets(ctx context.Context, obj *Bucket) ([]*Asset, error) {
	cmd := appbucket.GetBucketAssetsCommand{
		BucketID: obj.ID,
		Limit:    nil,
		LastKey:  nil,
	}

	assetIDs, err := r.bucketAppService.GetBucketAssets(ctx, cmd)
	if err != nil {
		return nil, err
	}

	result := make([]*Asset, 0, len(assetIDs))
	for _, assetID := range assetIDs {
		asset, err := r.assetAppService.GetAsset(ctx, appasset.GetAssetQuery{ID: assetID})
		if err != nil || asset == nil {
			continue
		}
		gqlAsset := domainAssetToGraphQL(asset)
		if gqlAsset == nil {
			continue
		}
		result = append(result, gqlAsset)
	}

	return result, nil
}

// CreateAsset is the resolver for the createAsset field.
func (r *mutationResolver) CreateAsset(ctx context.Context, input CreateAssetInput) (*Asset, error) {
	cmd := appasset.CreateAssetCommand{
		Slug:        input.Slug,
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Genre:       input.Genre,
		Genres:      input.Genres,
		Tags:        input.Tags,
		OwnerID:     input.OwnerID,
		ParentID:    input.ParentID,
		Metadata:    parseMetadata(input.Metadata),
	}

	asset, err := r.assetAppService.CreateAsset(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainAssetToGraphQL(asset), nil
}

// PatchAsset is the resolver for the patchAsset field.
func (r *mutationResolver) PatchAsset(ctx context.Context, id string, patches []*JSONPatch) (*Asset, error) {
	cmd := appasset.PatchAssetCommand{
		ID: id,
	}

	for _, patch := range patches {
		value := ""
		if patch.Value != nil {
			value = *patch.Value
		}
		cmd.Patches = append(cmd.Patches, appasset.JSONPatchOperation{
			Op:    patch.Op,
			Path:  patch.Path,
			Value: value,
		})
	}

	asset, err := r.assetAppService.PatchAsset(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainAssetToGraphQL(asset), nil
}

// DeleteAsset is the resolver for the deleteAsset field.
func (r *mutationResolver) DeleteAsset(ctx context.Context, id string) (bool, error) {
	cmd := appasset.DeleteAssetCommand{
		ID: id,
	}

	err := r.assetAppService.DeleteAsset(ctx, cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddVideo is the resolver for the addVideo field.
func (r *mutationResolver) AddVideo(ctx context.Context, input AddVideoInput) (*Video, error) {
	storageLocation, err := asset.NewS3Object(input.Bucket, input.Key, input.URL)
	if err != nil {
		return nil, err
	}

	cmd := appasset.AddVideoCommand{
		AssetID:         input.AssetID,
		Label:           input.Label,
		Format:          string(input.Format),
		StorageLocation: *storageLocation,
	}

	video, err := r.assetAppService.AddVideo(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainVideoToGraphQL(video), nil
}

// DeleteVideo is the resolver for the deleteVideo field.
func (r *mutationResolver) DeleteVideo(ctx context.Context, assetID string, videoID string) (*Asset, error) {
	cmd := appasset.RemoveVideoCommand{
		AssetID: assetID,
		VideoID: videoID,
	}
	err := r.assetAppService.RemoveVideo(ctx, cmd)
	if err != nil {
		return nil, err
	}
	asset, err := r.assetAppService.GetAsset(ctx, appasset.GetAssetQuery{ID: assetID})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(asset), nil
}

// CreateBucket is the resolver for the createBucket field.
func (r *mutationResolver) CreateBucket(ctx context.Context, input CreateBucketInput) (*Bucket, error) {
	cmd := appbucket.CreateBucketCommand{
		Key:         input.Key,
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     input.OwnerID,
		Status:      input.Status,
		Type:        &input.Type,
	}

	bucket, err := r.bucketAppService.CreateBucket(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainBucketToGraphQL(bucket), nil
}

// UpdateBucket is the resolver for the updateBucket field.
func (r *mutationResolver) UpdateBucket(ctx context.Context, input UpdateBucketInput) (*Bucket, error) {
	cmd := appbucket.UpdateBucketCommand{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     input.OwnerID,
		Metadata:    parseMetadata(input.Metadata),
		Status:      input.Status,
	}

	bucket, err := r.bucketAppService.UpdateBucket(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainBucketToGraphQL(bucket), nil
}

// DeleteBucket is the resolver for the deleteBucket field.
func (r *mutationResolver) DeleteBucket(ctx context.Context, input DeleteBucketInput) (bool, error) {
	cmd := appbucket.DeleteBucketCommand{
		ID:      input.ID,
		OwnerID: input.OwnerID,
	}

	err := r.bucketAppService.DeleteBucket(ctx, cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddAssetToBucket is the resolver for the addAssetToBucket field.
func (r *mutationResolver) AddAssetToBucket(ctx context.Context, input AddAssetToBucketInput) (bool, error) {
	cmd := appbucket.AddAssetToBucketCommand{
		BucketID: input.BucketID,
		AssetID:  input.AssetID,
		OwnerID:  input.OwnerID,
	}

	err := r.bucketAppService.AddAssetToBucket(ctx, cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RemoveAssetFromBucket is the resolver for the removeAssetFromBucket field.
func (r *mutationResolver) RemoveAssetFromBucket(ctx context.Context, input RemoveAssetFromBucketInput) (bool, error) {
	cmd := appbucket.RemoveAssetFromBucketCommand{
		BucketID: input.BucketID,
		AssetID:  input.AssetID,
		OwnerID:  input.OwnerID,
	}

	err := r.bucketAppService.RemoveAssetFromBucket(ctx, cmd)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddImage is the resolver for the addImage field.
func (r *mutationResolver) AddImage(ctx context.Context, input AddImageInput) (*Asset, error) {
	storageLocation, err := asset.NewS3Object(input.Bucket, input.Key, input.URL)
	if err != nil {
		return nil, err
	}

	domainType, err := asset.NewImageType(string(input.Type))
	if err != nil {
		return nil, err
	}

	domainImage, err := asset.NewImage(
		input.FileName,
		input.URL,
		domainType,
		storageLocation,
	)
	if err != nil {
		return nil, err
	}

	if err := domainImage.SetContentType(input.ContentType); err != nil {
		return nil, err
	}
	if err := domainImage.SetSize(int64(input.Size)); err != nil {
		return nil, err
	}

	cmd := appasset.AddImageCommand{
		AssetID: input.AssetID,
		Image:   *domainImage,
	}

	if err := r.assetAppService.AddImage(ctx, cmd); err != nil {
		return nil, err
	}

	assetObj, err := r.assetAppService.GetAsset(ctx, appasset.GetAssetQuery{ID: input.AssetID})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(assetObj), nil
}

// Assets is the resolver for the assets field.
func (r *queryResolver) Assets(ctx context.Context, limit *int, nextKey *string) (*AssetPage, error) {
	defer func() {
		if rec := recover(); rec != nil {
			logger.Get().Error(fmt.Sprintf("panic in Assets resolver: %v", rec))
		}
	}()

	limitInt := 0
	if limit != nil {
		limitInt = *limit
	}

	lastKeyMap := make(map[string]interface{})
	if nextKey != nil {
		lastKeyMap["key"] = *nextKey
	}

	cmd := appasset.ListAssetsQuery{
		Limit:   limitInt,
		LastKey: lastKeyMap,
	}

	page, err := r.assetAppService.ListAssets(ctx, cmd)
	if err != nil {
		logger.Get().Error(fmt.Sprintf("error in Assets resolver: %v", err))
		return nil, gqlerror.Errorf("internal system error")
	}

	result := make([]*Asset, len(page.Items))
	for i, asset := range page.Items {
		result[i] = domainAssetToGraphQL(asset)
	}

	var nextKeyStr *string
	if page.LastKey != nil {
		if key, ok := page.LastKey["key"].(string); ok {
			nextKeyStr = &key
		}
	}

	return &AssetPage{
		Items:   result,
		NextKey: nextKeyStr,
		HasMore: page.HasMore,
	}, nil
}

// Asset is the resolver for the asset field.
func (r *queryResolver) Asset(ctx context.Context, id *string) (*Asset, error) {
	if id == nil {
		return nil, nil
	}

	cmd := appasset.GetAssetQuery{
		ID: *id,
	}

	asset, err := r.assetAppService.GetAsset(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainAssetToGraphQL(asset), nil
}

// Buckets is the resolver for the buckets field.
func (r *queryResolver) Buckets(ctx context.Context, limit *int, nextKey *string) (*BucketPage, error) {
	defer func() {
		if rec := recover(); rec != nil {
			logger.Get().Error(fmt.Sprintf("panic in Buckets resolver: %v\n%s", rec, debug.Stack()))
		}
	}()

	lastKeyMap := make(map[string]interface{})
	if nextKey != nil {
		lastKeyMap["key"] = *nextKey
	}

	cmd := appbucket.ListBucketsCommand{
		Limit:   limit,
		LastKey: lastKeyMap,
	}

	page, err := r.bucketAppService.ListBuckets(ctx, cmd)
	if err != nil {
		logger.Get().Error(fmt.Sprintf("error in Buckets resolver: %v", err))
		return nil, gqlerror.Errorf("internal system error")
	}

	return domainBucketPageToGraphQL(page), nil
}

// Bucket is the resolver for the bucket field.
func (r *queryResolver) Bucket(ctx context.Context, id *string) (*Bucket, error) {
	if id == nil {
		return nil, nil
	}

	cmd := appbucket.GetBucketCommand{
		ID: *id,
	}

	bucket, err := r.bucketAppService.GetBucket(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainBucketToGraphQL(bucket), nil
}

// BucketByKey is the resolver for the bucketByKey field.
func (r *queryResolver) BucketByKey(ctx context.Context, key string) (*Bucket, error) {
	cmd := appbucket.GetBucketByKeyCommand{
		Key: key,
	}

	bucket, err := r.bucketAppService.GetBucketByKey(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainBucketToGraphQL(bucket), nil
}

// BucketsByOwner is the resolver for the bucketsByOwner field.
func (r *queryResolver) BucketsByOwner(ctx context.Context, ownerID string, limit *int, nextKey *string) (*BucketPage, error) {
	lastKeyMap := make(map[string]interface{})
	if nextKey != nil {
		lastKeyMap["key"] = *nextKey
	}

	cmd := appbucket.GetBucketsByOwnerCommand{
		OwnerID: ownerID,
		Limit:   limit,
		LastKey: lastKeyMap,
	}

	page, err := r.bucketAppService.GetBucketsByOwner(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainBucketPageToGraphQL(page), nil
}

// SearchBuckets is the resolver for the searchBuckets field.
func (r *queryResolver) SearchBuckets(ctx context.Context, query string, limit *int, nextKey *string) (*BucketPage, error) {
	lastKeyMap := make(map[string]interface{})
	if nextKey != nil {
		lastKeyMap["key"] = *nextKey
	}

	cmd := appbucket.SearchBucketsCommand{
		Query:   query,
		Limit:   limit,
		LastKey: lastKeyMap,
	}

	page, err := r.bucketAppService.SearchBuckets(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainBucketPageToGraphQL(page), nil
}

// SearchAssets is the resolver for the searchAssets field.
func (r *queryResolver) SearchAssets(ctx context.Context, query string, limit *int, nextKey *string) (*AssetPage, error) {
	lastKeyMap := make(map[string]interface{})
	if nextKey != nil {
		lastKeyMap["key"] = *nextKey
	}

	cmd := appasset.SearchAssetsQuery{
		Query:   query,
		Limit:   10,
		LastKey: lastKeyMap,
	}
	if limit != nil {
		cmd.Limit = *limit
	}

	page, err := r.assetAppService.SearchAssets(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return domainAssetPageToGraphQL(page), nil
}

// Bucket returns BucketResolver implementation.
func (r *Resolver) Bucket() BucketResolver { return &bucketResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type bucketResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
