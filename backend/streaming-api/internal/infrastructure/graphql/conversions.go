package graphql

import (
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

func normalizeGenreName(s string) string {
	if s == "" {
		return ""
	}
	v := strings.TrimSpace(strings.ToLower(s))
	v = strings.ReplaceAll(v, "-", "_")
	v = strings.ReplaceAll(v, " ", "_")
	return v
}

func ConvertGraphQLAssetToDomain(graphQLAsset *GraphQLAsset) (*entity.Asset, error) {
	assetID, err := assetvalueobjects.NewAssetID(graphQLAsset.ID)
	if err != nil {
		return nil, err
	}

	slug, err := assetvalueobjects.NewSlug(graphQLAsset.Slug)
	if err != nil {
		return nil, err
	}

	var title *assetvalueobjects.Title
	if graphQLAsset.Title != nil {
		titleVO, err := assetvalueobjects.NewTitle(*graphQLAsset.Title)
		if err != nil {
			return nil, err
		}
		title = titleVO
	}

	var description *assetvalueobjects.Description
	if graphQLAsset.Description != nil {
		descVO, err := assetvalueobjects.NewDescription(*graphQLAsset.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	assetType, err := assetvalueobjects.NewAssetType(graphQLAsset.Type)
	if err != nil {
		return nil, err
	}

	var genre *assetvalueobjects.Genre
	if graphQLAsset.Genre != nil {
		g := normalizeGenreName(*graphQLAsset.Genre)
		if g != "" {
			if genreVO, err := assetvalueobjects.NewGenre(g); err == nil {
				genre = genreVO
			}
		}
	}

	var genres *assetvalueobjects.Genres
	if len(graphQLAsset.Genres) > 0 {
		normalized := make([]string, 0, len(graphQLAsset.Genres))
		for _, g := range graphQLAsset.Genres {
			gn := normalizeGenreName(g)
			if gn != "" {
				if _, err := assetvalueobjects.NewGenre(gn); err == nil {
					normalized = append(normalized, gn)
				}
			}
		}
		if len(normalized) > 0 {
			if genresVO, err := assetvalueobjects.NewGenres(normalized); err == nil {
				genres = genresVO
			}
		}
	}

	var tags *assetvalueobjects.Tags
	if len(graphQLAsset.Tags) > 0 {
		tagsVO, err := assetvalueobjects.NewTags(graphQLAsset.Tags)
		if err != nil {
			return nil, err
		}
		tags = tagsVO
	}

	var status *assetvalueobjects.Status
	if graphQLAsset.Status != "" {
		statusVO, err := assetvalueobjects.NewStatus(graphQLAsset.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := assetvalueobjects.NewCreatedAt(graphQLAsset.CreatedAt)
	updatedAt := assetvalueobjects.NewUpdatedAt(graphQLAsset.UpdatedAt)

	var ownerID *assetvalueobjects.OwnerID
	if graphQLAsset.OwnerID != nil {
		ownerIDVO, err := assetvalueobjects.NewOwnerID(*graphQLAsset.OwnerID)
		if err != nil {
			return nil, err
		}
		ownerID = ownerIDVO
	}

	videos, err := ConvertGraphQLVideosToDomain(graphQLAsset.Videos)
	if err != nil {
		return nil, err
	}

	images, err := ConvertGraphQLImagesToDomain(graphQLAsset.Images)
	if err != nil {
		return nil, err
	}

	var publishRule *assetvalueobjects.PublishRuleValue
	if graphQLAsset.PublishRule != nil {
		publishRuleVO, err := assetvalueobjects.NewPublishRuleValue(
			graphQLAsset.PublishRule.PublishAt,
			graphQLAsset.PublishRule.UnpublishAt,
			graphQLAsset.PublishRule.Regions,
			graphQLAsset.PublishRule.AgeRating,
		)
		if err != nil {
			return nil, err
		}
		publishRule = publishRuleVO
	}

	return entity.NewAsset(
		*assetID,
		*slug,
		title,
		description,
		*assetType,
		genre,
		genres,
		tags,
		status,
		createdAt.Value(),
		updatedAt.Value(),
		graphQLAsset.Metadata,
		ownerID,
		videos,
		images,
		publishRule,
	), nil
}

func ConvertGraphQLVideosToDomain(graphQLVideos []GraphQLVideo) ([]entity.Video, error) {
	videos := make([]entity.Video, len(graphQLVideos))
	for i, graphQLVideo := range graphQLVideos {
		video, err := ConvertGraphQLVideoToDomain(graphQLVideo)
		if err != nil {
			return nil, err
		}
		videos[i] = *video
	}
	return videos, nil
}

func ConvertGraphQLVideoToDomain(graphQLVideo GraphQLVideo) (*entity.Video, error) {
	videoID, err := assetvalueobjects.NewVideoID(graphQLVideo.ID)
	if err != nil {
		return nil, err
	}

	var videoType *assetvalueobjects.VideoType
	if graphQLVideo.Type != "" {
		videoTypeVO, err := assetvalueobjects.NewVideoType(string(graphQLVideo.Type))
		if err != nil {
			return nil, err
		}
		videoType = videoTypeVO
	}

	var format *assetvalueobjects.VideoFormat
	if graphQLVideo.Format != "" {
		formatVO, err := assetvalueobjects.NewVideoFormat(string(graphQLVideo.Format))
		if err != nil {
			return nil, err
		}
		format = formatVO
	}

	storageLocation, err := assetvalueobjects.NewS3ObjectValue(
		graphQLVideo.StorageLocation.Bucket,
		graphQLVideo.StorageLocation.Key,
		graphQLVideo.StorageLocation.URL,
	)
	if err != nil {
		return nil, err
	}

	var streamInfo *assetvalueobjects.StreamInfoValue
	if graphQLVideo.StreamInfo != nil {
		streamInfoVO, err := assetvalueobjects.NewStreamInfoValue(
			graphQLVideo.StreamInfo.DownloadURL,
			graphQLVideo.StreamInfo.CDNPrefix,
			graphQLVideo.StreamInfo.URL,
		)
		if err != nil {
			return nil, err
		}
		streamInfo = streamInfoVO
	}

	var thumbnail *entity.Image
	if graphQLVideo.Thumbnail != nil {
		thumbnailVO, err := ConvertGraphQLImageToDomain(graphQLVideo.Thumbnail)
		if err != nil {
			return nil, err
		}
		thumbnail = thumbnailVO
	}

	metadata := ConvertStringSliceToString(graphQLVideo.Metadata)
	statusVO, err := assetvalueobjects.NewVideoStatus(string(graphQLVideo.Status))
	if err != nil {
		return nil, err
	}
	status := statusVO

	var quality *assetvalueobjects.VideoQuality
	if graphQLVideo.Quality != nil {
		qualityVO, err := assetvalueobjects.NewVideoQuality(*graphQLVideo.Quality)
		if err == nil {
			quality = qualityVO
		}
	}

	var transcodingInfo *assetvalueobjects.TranscodingInfo
	if graphQLVideo.TranscodingInfo != nil {
		transcodingInfo = &assetvalueobjects.TranscodingInfo{
			JobID:       graphQLVideo.TranscodingInfo.JobID,
			Progress:    graphQLVideo.TranscodingInfo.Progress,
			OutputURL:   graphQLVideo.TranscodingInfo.OutputURL,
			Error:       graphQLVideo.TranscodingInfo.Error,
			CompletedAt: graphQLVideo.TranscodingInfo.CompletedAt,
		}
	}

	return entity.NewVideo(
		*videoID,
		videoType,
		format,
		*storageLocation,
		graphQLVideo.Width,
		graphQLVideo.Height,
		graphQLVideo.Duration,
		graphQLVideo.Bitrate,
		graphQLVideo.Codec,
		graphQLVideo.Size,
		graphQLVideo.ContentType,
		streamInfo,
		metadata,
		status,
		thumbnail,
		graphQLVideo.CreatedAt,
		graphQLVideo.UpdatedAt,
		quality,
		graphQLVideo.IsReady,
		graphQLVideo.IsProcessing,
		graphQLVideo.IsFailed,
		graphQLVideo.SegmentCount,
		graphQLVideo.VideoCodec,
		graphQLVideo.AudioCodec,
		graphQLVideo.AvgSegmentDuration,
		graphQLVideo.Segments,
		graphQLVideo.FrameRate,
		graphQLVideo.AudioChannels,
		graphQLVideo.AudioSampleRate,
		transcodingInfo,
	), nil
}

func ConvertGraphQLImagesToDomain(graphQLImages []GraphQLImage) ([]entity.Image, error) {
	images := make([]entity.Image, len(graphQLImages))
	for i, graphQLImage := range graphQLImages {
		image, err := ConvertGraphQLImageToDomain(&graphQLImage)
		if err != nil {
			return nil, err
		}
		images[i] = *image
	}
	return images, nil
}

func ConvertGraphQLImageToDomain(graphQLImage *GraphQLImage) (*entity.Image, error) {
	imageID, err := assetvalueobjects.NewImageID(graphQLImage.ID)
	if err != nil {
		return nil, err
	}

	fileName, err := assetvalueobjects.NewFileName(graphQLImage.FileName)
	if err != nil {
		return nil, err
	}

	var imageType *assetvalueobjects.ImageType
	if graphQLImage.Type != "" {
		imageTypeVO, err := assetvalueobjects.NewImageType(string(graphQLImage.Type))
		if err != nil {
			return nil, err
		}
		imageType = imageTypeVO
	}

	var storageLocation *assetvalueobjects.S3ObjectValue
	if graphQLImage.StorageLocation != nil {
		storageLocationVO, err := assetvalueobjects.NewS3ObjectValue(
			graphQLImage.StorageLocation.Bucket,
			graphQLImage.StorageLocation.Key,
			graphQLImage.StorageLocation.URL,
		)
		if err != nil {
			return nil, err
		}
		storageLocation = storageLocationVO
	}

	metadata := ConvertStringSliceToString(graphQLImage.Metadata)

	return entity.NewImage(
		*imageID,
		*fileName,
		graphQLImage.URL,
		imageType,
		storageLocation,
		graphQLImage.Width,
		graphQLImage.Height,
		graphQLImage.Size,
		graphQLImage.ContentType,
		nil,
		metadata,
		graphQLImage.CreatedAt,
		graphQLImage.UpdatedAt,
	), nil
}

func ConvertGraphQLAssetsToDomain(graphQLAssets []GraphQLBucketAsset) ([]*entity.Asset, error) {
	assets := make([]*entity.Asset, len(graphQLAssets))
	for i, graphQLAsset := range graphQLAssets {
		domainAsset, err := ConvertGraphQLBucketAssetToDomain(graphQLAsset)
		if err != nil {
			return nil, err
		}
		assets[i] = domainAsset
	}
	return assets, nil
}

func ConvertGraphQLBucketAssetToDomain(graphQLAsset GraphQLBucketAsset) (*entity.Asset, error) {
	assetID, err := assetvalueobjects.NewAssetID(graphQLAsset.ID)
	if err != nil {
		return nil, err
	}

	slug, err := assetvalueobjects.NewSlug(graphQLAsset.Slug)
	if err != nil {
		return nil, err
	}

	var title *assetvalueobjects.Title
	if graphQLAsset.Title != nil {
		titleVO, err := assetvalueobjects.NewTitle(*graphQLAsset.Title)
		if err != nil {
			return nil, err
		}
		title = titleVO
	}

	var description *assetvalueobjects.Description
	if graphQLAsset.Description != nil {
		descVO, err := assetvalueobjects.NewDescription(*graphQLAsset.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	assetType, err := assetvalueobjects.NewAssetType(graphQLAsset.Type)
	if err != nil {
		return nil, err
	}

	var genre *assetvalueobjects.Genre
	if graphQLAsset.Genre != nil {
		g := normalizeGenreName(*graphQLAsset.Genre)
		if g != "" {
			if genreVO, err := assetvalueobjects.NewGenre(g); err == nil {
				genre = genreVO
			}
		}
	}

	var genres *assetvalueobjects.Genres
	if len(graphQLAsset.Genres) > 0 {
		normalized := make([]string, 0, len(graphQLAsset.Genres))
		for _, g := range graphQLAsset.Genres {
			gn := normalizeGenreName(g)
			if gn != "" {
				if _, err := assetvalueobjects.NewGenre(gn); err == nil {
					normalized = append(normalized, gn)
				}
			}
		}
		if len(normalized) > 0 {
			if genresVO, err := assetvalueobjects.NewGenres(normalized); err == nil {
				genres = genresVO
			}
		}
	}

	var tags *assetvalueobjects.Tags
	if len(graphQLAsset.Tags) > 0 {
		tagsVO, err := assetvalueobjects.NewTags(graphQLAsset.Tags)
		if err != nil {
			return nil, err
		}
		tags = tagsVO
	}

	var status *assetvalueobjects.Status
	if graphQLAsset.Status != "" {
		statusVO, err := assetvalueobjects.NewStatus(graphQLAsset.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := assetvalueobjects.NewCreatedAt(graphQLAsset.CreatedAt)
	updatedAt := assetvalueobjects.NewUpdatedAt(graphQLAsset.UpdatedAt)

	var ownerID *assetvalueobjects.OwnerID
	if graphQLAsset.OwnerID != nil {
		ownerIDVO, err := assetvalueobjects.NewOwnerID(*graphQLAsset.OwnerID)
		if err != nil {
			return nil, err
		}
		ownerID = ownerIDVO
	}

	videos, err := ConvertGraphQLVideosToDomain(graphQLAsset.Videos)
	if err != nil {
		return nil, err
	}

	images, err := ConvertGraphQLImagesToDomain(graphQLAsset.Images)
	if err != nil {
		return nil, err
	}

	var publishRule *assetvalueobjects.PublishRuleValue
	if graphQLAsset.PublishRule != nil {
		publishRuleVO, err := assetvalueobjects.NewPublishRuleValue(
			graphQLAsset.PublishRule.PublishAt,
			graphQLAsset.PublishRule.UnpublishAt,
			graphQLAsset.PublishRule.Regions,
			graphQLAsset.PublishRule.AgeRating,
		)
		if err != nil {
			return nil, err
		}
		publishRule = publishRuleVO
	}

	return entity.NewAsset(
		*assetID,
		*slug,
		title,
		description,
		*assetType,
		genre,
		genres,
		tags,
		status,
		createdAt.Value(),
		updatedAt.Value(),
		graphQLAsset.Metadata,
		ownerID,
		videos,
		images,
		publishRule,
	), nil
}

func ConvertStringSliceToString(slice []string) *string {
	if len(slice) == 0 {
		return nil
	}
	// TODO: Consider JSON marshaling for complex metadata
	result := strings.Join(slice, ",")
	return &result
}
