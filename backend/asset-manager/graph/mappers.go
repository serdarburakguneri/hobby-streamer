package graph

import (
	"encoding/json"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/graph/model"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func mapAssetToGraphQL(a *asset.Asset) *model.Asset {
	if a == nil {
		return nil
	}

	status := a.Status()
	graphQLAsset := &model.Asset{
		ID:          a.ID,
		Slug:        a.Slug,
		Title:       a.Title,
		Description: a.Description,
		Genre:       a.Genre,
		Genres:      a.Genres,
		Tags:        a.Tags,
		Status:      &status,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		OwnerID:     a.OwnerID,
		Videos:      mapVideosToGraphQL(a.Videos),
		Images:      mapImagesToGraphQL(a.Images),
		PublishRule: mapPublishRuleToGraphQL(a.PublishRule),
	}

	// Map parent information
	if a.Parent != nil {
		graphQLAsset.Parent = mapAssetToGraphQL(a.Parent)
	}

	if a.Type != nil {
		switch *a.Type {
		case asset.AssetTypeMovie:
			graphQLAsset.Type = model.AssetTypeMovie
		case asset.AssetTypeSeries:
			graphQLAsset.Type = model.AssetTypeSeries
		case asset.AssetTypeSeason:
			graphQLAsset.Type = model.AssetTypeSeason
		case asset.AssetTypeEpisode:
			graphQLAsset.Type = model.AssetTypeEpisode
		case asset.AssetTypeDocumentary:
			graphQLAsset.Type = model.AssetTypeDocumentary
		}
	}

	if a.Metadata != nil {
		metadataJSON, err := json.Marshal(a.Metadata)
		if err == nil {
			metadataStr := string(metadataJSON)
			graphQLAsset.Metadata = &metadataStr
		}
	}

	return graphQLAsset
}

func mapAssetsToGraphQL(assets []asset.Asset) []*model.Asset {
	if assets == nil {
		return nil
	}

	result := make([]*model.Asset, len(assets))
	for i, a := range assets {
		result[i] = mapAssetToGraphQL(&a)
	}
	return result
}

func mapGraphQLAssetInputToAsset(input model.AssetInput) *asset.Asset {
	log := logger.Get().WithService("graphql-mapper")

	assetModel := &asset.Asset{
		ID:          "",
		Slug:        input.Slug,
		Title:       input.Title,
		Description: input.Description,
		Genre:       input.Genre,
		Genres:      input.Genres,
		Tags:        input.Tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     input.OwnerID,
		ParentID:    input.ParentID,
	}

	switch input.Type {
	case model.AssetTypeMovie:
		assetType := asset.AssetTypeMovie
		assetModel.Type = &assetType
	case model.AssetTypeSeries:
		assetType := asset.AssetTypeSeries
		assetModel.Type = &assetType
	case model.AssetTypeSeason:
		assetType := asset.AssetTypeSeason
		assetModel.Type = &assetType
	case model.AssetTypeEpisode:
		assetType := asset.AssetTypeEpisode
		assetModel.Type = &assetType
	case model.AssetTypeDocumentary:
		assetType := asset.AssetTypeDocumentary
		assetModel.Type = &assetType
	case model.AssetTypeMusic:
		assetType := asset.AssetTypeMusic
		assetModel.Type = &assetType
	case model.AssetTypePodcast:
		assetType := asset.AssetTypePodcast
		assetModel.Type = &assetType
	case model.AssetTypeTrailer:
		assetType := asset.AssetTypeTrailer
		assetModel.Type = &assetType
	case model.AssetTypeBehindTheScenes:
		assetType := asset.AssetTypeBehindTheScenes
		assetModel.Type = &assetType
	case model.AssetTypeInterview:
		assetType := asset.AssetTypeInterview
		assetModel.Type = &assetType
	}

	if input.Metadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(*input.Metadata), &metadata); err == nil {
			assetModel.Metadata = metadata
		} else {
			log.WithError(err).Error("Failed to parse metadata JSON")
		}
	} else {
		assetModel.Metadata = make(map[string]interface{})
	}

	return assetModel
}

func mapBucketToGraphQL(b *bucket.Bucket) *model.Bucket {
	if b == nil {
		return nil
	}

	graphQLBucket := &model.Bucket{
		ID:          b.ID,
		Key:         b.Key,
		Name:        b.Name,
		Description: &b.Description,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
		AssetIds:    b.AssetIDs,
	}

	switch b.Type {
	case bucket.BucketTypePlaylist:
		graphQLBucket.Type = model.BucketTypePlaylist
	case bucket.BucketTypeCollection:
		graphQLBucket.Type = model.BucketTypeCollection
	case bucket.BucketTypeCategory:
		graphQLBucket.Type = model.BucketTypeCategory
	}

	switch b.Status {
	case bucket.BucketStatusActive:
		status := model.BucketStatusActive
		graphQLBucket.Status = &status
	case bucket.BucketStatusInactive:
		status := model.BucketStatusInactive
		graphQLBucket.Status = &status
	case bucket.BucketStatusDraft:
		status := model.BucketStatusDraft
		graphQLBucket.Status = &status
	}

	return graphQLBucket
}

func mapBucketsToGraphQL(buckets []bucket.Bucket) []*model.Bucket {
	if buckets == nil {
		return nil
	}

	result := make([]*model.Bucket, len(buckets))
	for i, b := range buckets {
		result[i] = mapBucketToGraphQL(&b)
	}
	return result
}

func mapGraphQLBucketInputToBucket(input model.BucketInput) *bucket.Bucket {
	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	bucketModel := &bucket.Bucket{
		ID:          "",
		Key:         input.Key,
		Name:        input.Name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AssetIDs:    input.AssetIds,
	}

	switch input.Type {
	case model.BucketTypePlaylist:
		bucketModel.Type = bucket.BucketTypePlaylist
	case model.BucketTypeCollection:
		bucketModel.Type = bucket.BucketTypeCollection
	case model.BucketTypeCategory:
		bucketModel.Type = bucket.BucketTypeCategory
	}

	if input.Status != nil {
		switch *input.Status {
		case model.BucketStatusActive:
			bucketModel.Status = bucket.BucketStatusActive
		case model.BucketStatusInactive:
			bucketModel.Status = bucket.BucketStatusInactive
		case model.BucketStatusDraft:
			bucketModel.Status = bucket.BucketStatusDraft
		}
	}

	return bucketModel
}

func mapGraphQLUpdateBucketInputToBucket(input model.UpdateBucketInput) *bucket.Bucket {
	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	bucketModel := &bucket.Bucket{
		Name:        input.Name,
		Description: description,
	}

	switch input.Type {
	case model.BucketTypePlaylist:
		bucketModel.Type = bucket.BucketTypePlaylist
	case model.BucketTypeCollection:
		bucketModel.Type = bucket.BucketTypeCollection
	case model.BucketTypeCategory:
		bucketModel.Type = bucket.BucketTypeCategory
	}

	if input.Status != nil {
		switch *input.Status {
		case model.BucketStatusActive:
			bucketModel.Status = bucket.BucketStatusActive
		case model.BucketStatusInactive:
			bucketModel.Status = bucket.BucketStatusInactive
		case model.BucketStatusDraft:
			bucketModel.Status = bucket.BucketStatusDraft
		}
	}

	return bucketModel
}

// Pagination mappers

func mapAssetPageToGraphQL(page *asset.AssetPage) *model.AssetPage {
	if page == nil {
		return nil
	}

	return &model.AssetPage{
		Items: mapAssetsToGraphQL(page.Items),
	}
}

func mapBucketPageToGraphQL(page *bucket.BucketPage) *model.BucketPage {
	if page == nil {
		return nil
	}

	return &model.BucketPage{
		Items: mapBucketsToGraphQL(page.Items),
	}
}

func mapPublishRuleToGraphQL(rule *asset.PublishRule) *model.PublishRule {
	if rule == nil {
		return nil
	}

	return &model.PublishRule{
		PublishAt:   &rule.PublishAt,
		UnpublishAt: &rule.UnpublishAt,
		Regions:     rule.Regions,
		AgeRating:   &rule.AgeRating,
	}
}

func mapVideosToGraphQL(videos []asset.Video) []*model.Video {
	if videos == nil {
		return nil
	}

	result := make([]*model.Video, 0, len(videos))
	for _, video := range videos {
		var videoType model.VideoType
		switch video.Type {
		case asset.VideoTypeMain:
			videoType = model.VideoTypeMain
		case asset.VideoTypeTrailer:
			videoType = model.VideoTypeTrailer
		case asset.VideoTypeBehind:
			videoType = model.VideoTypeBehindTheScenes
		case asset.VideoTypeInterview:
			videoType = model.VideoTypeInterview
		}

		graphQLVideo := &model.Video{
			ID:              video.ID,
			Type:            videoType,
			Format:          string(video.Format),
			StorageLocation: mapS3ObjectToGraphQL(&video.StorageLocation),
			Width:           &video.Width,
			Height:          &video.Height,
			Duration:        &video.Duration,
			Bitrate:         &video.Bitrate,
			Codec:           &video.Codec,
			Size:            &[]int{int(video.Size)}[0],
			ContentType:     &video.ContentType,
			StreamInfo:      mapStreamInfoToGraphQL(video.StreamInfo),
			Status:          &video.Status,
			Thumbnail:       mapImageToGraphQL(video.Thumbnail),
			CreatedAt:       video.CreatedAt,
			UpdatedAt:       video.UpdatedAt,
		}

		if video.Metadata != nil {
			metadataJSON, err := json.Marshal(video.Metadata)
			if err == nil {
				metadataStr := string(metadataJSON)
				graphQLVideo.Metadata = &metadataStr
			}
		}

		result = append(result, graphQLVideo)
	}
	return result
}

func mapImagesToGraphQL(images []asset.Image) []*model.Image {
	if images == nil {
		return nil
	}

	result := make([]*model.Image, 0, len(images))
	for _, img := range images {
		result = append(result, mapImageToGraphQL(&img))
	}
	return result
}

func mapImageToGraphQL(img *asset.Image) *model.Image {
	if img == nil {
		return nil
	}

	var imageType model.ImageType
	switch img.Type {
	case asset.ImageTypeThumbnail:
		imageType = model.ImageTypeThumbnail
	case asset.ImageTypePoster:
		imageType = model.ImageTypePoster
	case asset.ImageTypeBanner:
		imageType = model.ImageTypeBanner
	case asset.ImageTypeHero:
		imageType = model.ImageTypeHero
	case asset.ImageTypeLogo:
		imageType = model.ImageTypeLogo
	case asset.ImageTypeScreenshot:
		imageType = model.ImageTypeScreenshot

	default:
		logger.Get().WithService("graph-mapper").Error("Unknown image type, using POSTER as fallback", "image_type", img.Type, "image_id", img.ID, "fileName", img.FileName)
		imageType = model.ImageTypePoster
	}

	graphQLImage := &model.Image{
		ID:          img.ID,
		FileName:    img.FileName,
		URL:         img.URL,
		Type:        imageType,
		Width:       &img.Width,
		Height:      &img.Height,
		Size:        &[]int{int(img.Size)}[0],
		ContentType: &img.ContentType,
		CreatedAt:   img.CreatedAt,
		UpdatedAt:   img.UpdatedAt,
	}

	if img.StorageLocation != nil {
		graphQLImage.StorageLocation = mapS3ObjectToGraphQL(img.StorageLocation)
	}

	if img.Metadata != nil {
		metadataJSON, err := json.Marshal(img.Metadata)
		if err == nil {
			metadataStr := string(metadataJSON)
			graphQLImage.Metadata = &metadataStr
		}
	}

	return graphQLImage
}

func mapS3ObjectToGraphQL(obj *asset.S3Object) *model.S3Object {
	if obj == nil {
		return nil
	}

	return &model.S3Object{
		Bucket: obj.Bucket,
		Key:    obj.Key,
		URL:    obj.URL,
	}
}

func mapStreamInfoToGraphQL(info *asset.StreamInfo) *model.StreamInfo {
	if info == nil {
		return nil
	}

	return &model.StreamInfo{
		DownloadURL: info.DownloadURL,
		CdnPrefix:   info.CdnPrefix,
		PlayURL:     info.PlayURL,
	}
}
