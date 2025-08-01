package graphql

import (
	"context"
	"fmt"

	assetCommands "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	assetAppQueries "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/queries"
	assetvo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
)

// CreateAsset handles createAsset GraphQL mutation.
func (r *mutationResolver) CreateAsset(ctx context.Context, input CreateAssetInput) (*Asset, error) {
	cmd, err := MapCreateAssetInput(input)
	if err != nil {
		return nil, err
	}
	a, err := r.assetCommandService.CreateAsset(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

// PatchAsset handles patchAsset GraphQL mutation.
func (r *mutationResolver) PatchAsset(ctx context.Context, id string, patches []*JSONPatch) (*Asset, error) {
	cmd, err := MapPatchAssetInput(id, patches)
	if err != nil {
		return nil, err
	}
	if err := r.assetCommandService.PatchAsset(ctx, cmd); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: id})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

// DeleteAsset handles deleteAsset GraphQL mutation.
func (r *mutationResolver) DeleteAsset(ctx context.Context, id string) (bool, error) {
	cmd, err := MapDeleteAssetInput(id)
	if err != nil {
		return false, err
	}
	if err := r.assetCommandService.DeleteAsset(ctx, cmd); err != nil {
		return false, err
	}
	return true, nil
}

// AddVideo handles addVideo GraphQL mutation.
func (r *mutationResolver) AddVideo(ctx context.Context, input AddVideoInput) (*Video, error) {
	idVO, err := assetvo.NewAssetID(input.AssetID)
	if err != nil {
		return nil, err
	}
	formatVO, err := assetvo.NewVideoFormat(string(input.Format))
	if err != nil {
		return nil, err
	}
	s3VO, err := assetvo.NewS3Object(input.Bucket, input.Key, input.URL)
	if err != nil {
		return nil, err
	}
	cmd := assetCommands.AddVideoCommand{
		AssetID:         *idVO,
		Label:           input.Label,
		Format:          formatVO,
		StorageLocation: *s3VO,
		Size:            int64(input.Size),
		ContentType:     input.ContentType,
	}
	if err := r.assetCommandService.AddVideo(ctx, cmd); err != nil {
		return nil, err
	}
	assetEntity, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: input.AssetID})
	if err != nil {
		return nil, err
	}
	for _, v := range assetEntity.Videos() {
		if v.StorageLocation().Key() == input.Key {
			return domainVideoToGraphQL(v), nil
		}
	}
	return nil, fmt.Errorf("video not found after creation")
}

// DeleteVideo handles deleteVideo GraphQL mutation.
func (r *mutationResolver) DeleteVideo(ctx context.Context, assetID string, videoID string) (*Asset, error) {
	idVO, err := assetvo.NewAssetID(assetID)
	if err != nil {
		return nil, err
	}
	if err := r.assetCommandService.RemoveVideo(ctx, assetCommands.RemoveVideoCommand{AssetID: *idVO, VideoID: videoID}); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: assetID})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

// AddImage handles addImage GraphQL mutation.
func (r *mutationResolver) AddImage(ctx context.Context, input AddImageInput) (*Asset, error) {
	idVO, err := assetvo.NewAssetID(input.AssetID)
	if err != nil {
		return nil, err
	}
	ito, err := assetvo.NewImageType(string(input.Type))
	if err != nil {
		return nil, err
	}
	imgVO, err := assetvo.NewImage(input.FileName, input.URL, *ito, input.ContentType)
	if err != nil {
		return nil, err
	}
	if err := r.assetCommandService.AddImage(ctx, assetCommands.AddImageCommand{AssetID: *idVO, Image: *imgVO}); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: input.AssetID})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

// Assets lists assets.
func (r *queryResolver) Assets(ctx context.Context, limit *int, offset *int) ([]*Asset, error) {
	q := assetAppQueries.ListAssetsQuery{Limit: limit, Offset: offset}
	items, err := r.assetQueryService.ListAssets(ctx, q)
	if err != nil {
		return nil, err
	}
	out := make([]*Asset, len(items))
	for i, a := range items {
		out[i] = domainAssetToGraphQL(a)
	}
	return out, nil
}

// Asset retrieves a single asset.
func (r *queryResolver) Asset(ctx context.Context, id *string) (*Asset, error) {
	if id == nil {
		return nil, nil
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: *id})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

// SearchAssets paginates assets.
func (r *queryResolver) SearchAssets(ctx context.Context, query string, limit *int, offset *int) ([]*Asset, error) {
	page, err := r.assetQueryService.SearchAssetsPage(ctx, assetAppQueries.SearchAssetsQuery{Query: query, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*Asset, len(page.Items))
	for i, a := range page.Items {
		out[i] = domainAssetToGraphQL(a)
	}
	return out, nil
}
