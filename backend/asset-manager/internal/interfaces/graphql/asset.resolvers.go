package graphql

import (
	"context"
	"fmt"
	"time"

	assetCommands "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	assetAppQueries "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/queries"
	assetvo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
)

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
	upsert := assetCommands.UpsertVideoCommand{
		AssetID:         *idVO,
		Label:           input.Label,
		Format:          formatVO,
		StorageLocation: *s3VO,
		Size:            int64(input.Size),
		ContentType:     input.ContentType,
	}
	if _, _, err := r.assetCommandService.UpsertVideo(ctx, upsert); err != nil {
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

	cdnPrefix, playURL := r.cdnService.BuildPlayURL(input.Key)
	if si, err := assetvo.NewStreamInfo(nil, &cdnPrefix, &playURL); err == nil {
		imgVO.SetStreamInfo(si)
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

func (r *mutationResolver) UpdateAssetTitle(ctx context.Context, id string, title string) (*Asset, error) {
	idVO, err := assetvo.NewAssetID(id)
	if err != nil {
		return nil, err
	}
	titleVO, err := assetvo.NewTitle(title)
	if err != nil {
		return nil, err
	}
	cmd := assetCommands.UpdateAssetTitleCommand{AssetID: *idVO, Title: *titleVO}
	if err := r.assetCommandService.UpdateAssetTitle(ctx, cmd); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: id})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

func (r *mutationResolver) UpdateAssetDescription(ctx context.Context, id string, description string) (*Asset, error) {
	idVO, err := assetvo.NewAssetID(id)
	if err != nil {
		return nil, err
	}
	descVO, err := assetvo.NewDescription(description)
	if err != nil {
		return nil, err
	}
	cmd := assetCommands.UpdateAssetDescriptionCommand{AssetID: *idVO, Description: *descVO}
	if err := r.assetCommandService.UpdateAssetDescription(ctx, cmd); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: id})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

func (r *mutationResolver) SetAssetPublishRule(ctx context.Context, id string, rule PublishRuleInput) (*Asset, error) {
	idVO, err := assetvo.NewAssetID(id)
	if err != nil {
		return nil, err
	}
	var publishAtPtr *time.Time
	if rule.PublishAt != nil {
		t := *rule.PublishAt
		publishAtPtr = &t
	}
	var unpublishAtPtr *time.Time
	if rule.UnpublishAt != nil {
		t := *rule.UnpublishAt
		unpublishAtPtr = &t
	}
	regions := make([]string, len(rule.Regions))
	copy(regions, rule.Regions)
	var ageRatingPtr *string
	if rule.AgeRating != nil {
		a := *rule.AgeRating
		ageRatingPtr = &a
	}
	pr, err := assetvo.NewPublishRule(publishAtPtr, unpublishAtPtr, regions, ageRatingPtr)
	if err != nil {
		return nil, err
	}
	cmd := assetCommands.SetAssetPublishRuleCommand{AssetID: *idVO, PublishRule: *pr}
	if err := r.assetCommandService.SetPublishRule(ctx, cmd); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: id})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

func (r *mutationResolver) ClearAssetPublishRule(ctx context.Context, id string) (*Asset, error) {
	idVO, err := assetvo.NewAssetID(id)
	if err != nil {
		return nil, err
	}
	cmd := assetCommands.ClearAssetPublishRuleCommand{AssetID: *idVO}
	if err := r.assetCommandService.ClearPublishRule(ctx, cmd); err != nil {
		return nil, err
	}
	a, err := r.assetQueryService.GetAsset(ctx, assetAppQueries.GetAssetQuery{ID: id})
	if err != nil {
		return nil, err
	}
	return domainAssetToGraphQL(a), nil
}

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
