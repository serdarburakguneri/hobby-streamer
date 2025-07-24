package graphql

import (
	"context"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
)

type AssetRepository struct {
	client         *Client
	logger         *logger.Logger
	circuitBreaker *errors.CircuitBreaker
}

func NewAssetRepository(client *Client, circuitBreaker *errors.CircuitBreaker) asset.Repository {
	return &AssetRepository{
		client:         client,
		logger:         logger.Get().WithService("asset-graphql-repository"),
		circuitBreaker: circuitBreaker,
	}
}

func (r *AssetRepository) GetByID(ctx context.Context, id asset.AssetID) (*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range assets {
		if a.ID().Equals(id) {
			return a, nil
		}
	}

	return nil, errors.NewNotFoundError("asset not found", nil)
}

func (r *AssetRepository) GetBySlug(ctx context.Context, slug asset.Slug) (*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range assets {
		if a.Slug().Equals(slug) {
			return a, nil
		}
	}

	return nil, errors.NewNotFoundError("asset not found", nil)
}

func (r *AssetRepository) GetAll(ctx context.Context) ([]*asset.Asset, error) {
	query := `
		query GetAssets {
			assets {
				items {
					id
					slug
					title
					description
					type
					genre
					genres
					tags
					status
					createdAt
					updatedAt
					metadata
					ownerId
					videos {
						id
						label
						type
						format
						storageLocation {
							bucket
							key
							url
						}
						width
						height
						duration
						bitrate
						codec
						size
						contentType
						streamInfo {
							downloadUrl
							cdnPrefix
							url
						}
						metadata
						status
						thumbnail {
							id
							fileName
							url
							type
							storageLocation {
								bucket
								key
								url
							}
							width
							height
							size
							contentType
							metadata
							createdAt
							updatedAt
						}
						createdAt
						updatedAt
						quality
						isReady
						isProcessing
						isFailed
						segmentCount
						videoCodec
						audioCodec
						avgSegmentDuration
						segments
						frameRate
						audioChannels
						audioSampleRate
						transcodingInfo {
							jobId
							progress
							outputUrl
							error
							completedAt
						}
					}
					images {
						id
						fileName
						url
						type
						storageLocation {
							bucket
							key
							url
						}
						width
						height
						size
						contentType
						metadata
						createdAt
						updatedAt
					}
					publishRule {
						publishAt
						unpublishAt
						regions
						ageRating
					}
				}
			}
		}
	`

	var response struct {
		Assets struct {
			Items []*GraphQLAsset `json:"items"`
		} `json:"assets"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, nil, &response)
	})

	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_all_assets",
		})
	}

	assets := make([]*asset.Asset, len(response.Assets.Items))
	for i, graphQLAsset := range response.Assets.Items {
		domainAsset, err := r.convertGraphQLAssetToDomain(graphQLAsset)
		if err != nil {
			r.logger.WithError(err).Error("Failed to convert GraphQL asset to domain", "asset_id", graphQLAsset.ID)
			continue
		}
		assets[i] = domainAsset
	}

	return assets, nil
}

func (r *AssetRepository) GetPublic(ctx context.Context) ([]*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var publicAssets []*asset.Asset
	for _, a := range assets {
		if a.IsPublished() {
			publicAssets = append(publicAssets, a)
		}
	}

	return publicAssets, nil
}

func (r *AssetRepository) GetByType(ctx context.Context, assetType asset.AssetType) ([]*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*asset.Asset
	for _, a := range assets {
		if a.Type().Equals(assetType) {
			filteredAssets = append(filteredAssets, a)
		}
	}

	return filteredAssets, nil
}

func (r *AssetRepository) GetByGenre(ctx context.Context, genre asset.Genre) ([]*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*asset.Asset
	for _, a := range assets {
		if a.Genre() != nil && a.Genre().Equals(genre) {
			filteredAssets = append(filteredAssets, a)
		}
	}

	return filteredAssets, nil
}

func (r *AssetRepository) GetByOwner(ctx context.Context, ownerID asset.OwnerID) ([]*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*asset.Asset
	for _, a := range assets {
		if a.OwnerID() != nil && a.OwnerID().Equals(ownerID) {
			filteredAssets = append(filteredAssets, a)
		}
	}

	return filteredAssets, nil
}

func (r *AssetRepository) GetReady(ctx context.Context) ([]*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var readyAssets []*asset.Asset
	for _, a := range assets {
		if a.IsReady() {
			readyAssets = append(readyAssets, a)
		}
	}

	return readyAssets, nil
}

func (r *AssetRepository) GetPublished(ctx context.Context) ([]*asset.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var publishedAssets []*asset.Asset
	for _, a := range assets {
		if a.IsPublished() {
			publishedAssets = append(publishedAssets, a)
		}
	}

	return publishedAssets, nil
}

func (r *AssetRepository) Search(ctx context.Context, query string, filters *asset.SearchFilters) ([]*asset.Asset, error) {
	_, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	searchService := asset.NewAssetSearchService()
	return searchService.SearchAssets(query, filters)
}

func (r *AssetRepository) GetRecommended(ctx context.Context, asset *asset.Asset, limit int) ([]*asset.Asset, error) {
	return nil, errors.NewNotFoundError("get recommended assets not implemented", nil)
}

func (r *AssetRepository) GetStreamingInfo(ctx context.Context, asset *asset.Asset, userID string, region string, userAge int) (*asset.StreamingInfo, error) {
	return nil, errors.NewNotFoundError("get streaming info not implemented", nil)
}

func (r *AssetRepository) convertGraphQLAssetToDomain(graphQLAsset *GraphQLAsset) (*asset.Asset, error) {
	assetID, err := asset.NewAssetID(graphQLAsset.ID)
	if err != nil {
		return nil, err
	}

	slug, err := asset.NewSlug(graphQLAsset.Slug)
	if err != nil {
		return nil, err
	}

	var title *asset.Title
	if graphQLAsset.Title != nil {
		titleVO, err := asset.NewTitle(*graphQLAsset.Title)
		if err != nil {
			return nil, err
		}
		title = titleVO
	}

	var description *asset.Description
	if graphQLAsset.Description != nil {
		descVO, err := asset.NewDescription(*graphQLAsset.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	assetType, err := asset.NewAssetType(graphQLAsset.Type)
	if err != nil {
		return nil, err
	}

	var genre *asset.Genre
	if graphQLAsset.Genre != nil {
		genreVO, err := asset.NewGenre(*graphQLAsset.Genre)
		if err != nil {
			return nil, err
		}
		genre = genreVO
	}

	var genres *asset.Genres
	if len(graphQLAsset.Genres) > 0 {
		genresVO, err := asset.NewGenres(graphQLAsset.Genres)
		if err != nil {
			return nil, err
		}
		genres = genresVO
	}

	var tags *asset.Tags
	if len(graphQLAsset.Tags) > 0 {
		tagsVO, err := asset.NewTags(graphQLAsset.Tags)
		if err != nil {
			return nil, err
		}
		tags = tagsVO
	}

	var status *asset.Status
	if graphQLAsset.Status != "" {
		statusVO, err := asset.NewStatus(graphQLAsset.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := asset.NewCreatedAt(graphQLAsset.CreatedAt)
	updatedAt := asset.NewUpdatedAt(graphQLAsset.UpdatedAt)

	var ownerID *asset.OwnerID
	if graphQLAsset.OwnerID != nil {
		ownerIDVO, err := asset.NewOwnerID(*graphQLAsset.OwnerID)
		if err != nil {
			return nil, err
		}
		ownerID = ownerIDVO
	}

	videos, err := r.convertGraphQLVideosToDomain(graphQLAsset.Videos)
	if err != nil {
		return nil, err
	}

	images, err := r.convertGraphQLImagesToDomain(graphQLAsset.Images)
	if err != nil {
		return nil, err
	}

	var publishRule *asset.PublishRuleValue
	if graphQLAsset.PublishRule != nil {
		publishRuleVO, err := asset.NewPublishRuleValue(
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

	return asset.NewAsset(
		*assetID,
		*slug,
		title,
		description,
		*assetType,
		genre,
		genres,
		tags,
		status,
		*createdAt,
		*updatedAt,
		graphQLAsset.Metadata,
		ownerID,
		videos,
		images,
		publishRule,
	), nil
}

func (r *AssetRepository) convertGraphQLVideosToDomain(graphQLVideos []GraphQLVideo) ([]asset.Video, error) {
	videos := make([]asset.Video, len(graphQLVideos))
	for i, graphQLVideo := range graphQLVideos {
		video, err := r.convertGraphQLVideoToDomain(graphQLVideo)
		if err != nil {
			return nil, err
		}
		videos[i] = *video
	}
	return videos, nil
}

func (r *AssetRepository) convertGraphQLVideoToDomain(graphQLVideo GraphQLVideo) (*asset.Video, error) {
	videoID, err := asset.NewVideoID(graphQLVideo.ID)
	if err != nil {
		return nil, err
	}

	var videoType *asset.VideoType
	if graphQLVideo.Type != "" {
		videoTypeVO, err := asset.NewVideoType(string(graphQLVideo.Type))
		if err != nil {
			return nil, err
		}
		videoType = videoTypeVO
	}

	var format *asset.VideoFormat
	if graphQLVideo.Format != "" {
		formatVO, err := asset.NewVideoFormat(string(graphQLVideo.Format))
		if err != nil {
			return nil, err
		}
		format = formatVO
	}

	storageLocation, err := asset.NewS3ObjectValue(
		graphQLVideo.StorageLocation.Bucket,
		graphQLVideo.StorageLocation.Key,
		graphQLVideo.StorageLocation.URL,
	)
	if err != nil {
		return nil, err
	}

	var streamInfo *asset.StreamInfoValue
	if graphQLVideo.StreamInfo != nil {
		streamInfoVO, err := asset.NewStreamInfoValue(
			graphQLVideo.StreamInfo.DownloadURL,
			graphQLVideo.StreamInfo.CDNPrefix,
			graphQLVideo.StreamInfo.URL,
		)
		if err != nil {
			return nil, err
		}
		streamInfo = streamInfoVO
	}

	var thumbnail *asset.Image
	if graphQLVideo.Thumbnail != nil {
		thumbnailVO, err := r.convertGraphQLImageToDomain(graphQLVideo.Thumbnail)
		if err != nil {
			return nil, err
		}
		thumbnail = thumbnailVO
	}

	metadata := convertStringSliceToString(graphQLVideo.Metadata)
	status := string(graphQLVideo.Status)

	var quality *asset.VideoQuality
	if graphQLVideo.Quality != nil {
		qualityVO, err := asset.NewVideoQuality(*graphQLVideo.Quality)
		if err == nil {
			quality = qualityVO
		}
	}

	var transcodingInfo *asset.TranscodingInfo
	if graphQLVideo.TranscodingInfo != nil {
		transcodingInfo = &asset.TranscodingInfo{
			JobID:       graphQLVideo.TranscodingInfo.JobID,
			Progress:    graphQLVideo.TranscodingInfo.Progress,
			OutputURL:   graphQLVideo.TranscodingInfo.OutputURL,
			Error:       graphQLVideo.TranscodingInfo.Error,
			CompletedAt: graphQLVideo.TranscodingInfo.CompletedAt,
		}
	}

	return asset.NewVideo(
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
		&status,
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

func (r *AssetRepository) convertGraphQLImagesToDomain(graphQLImages []GraphQLImage) ([]asset.Image, error) {
	images := make([]asset.Image, len(graphQLImages))
	for i, graphQLImage := range graphQLImages {
		image, err := r.convertGraphQLImageToDomain(&graphQLImage)
		if err != nil {
			return nil, err
		}
		images[i] = *image
	}
	return images, nil
}

func (r *AssetRepository) convertGraphQLImageToDomain(graphQLImage *GraphQLImage) (*asset.Image, error) {
	imageID, err := asset.NewImageID(graphQLImage.ID)
	if err != nil {
		return nil, err
	}

	fileName, err := asset.NewFileName(graphQLImage.FileName)
	if err != nil {
		return nil, err
	}

	var imageType *asset.ImageType
	if graphQLImage.Type != "" {
		imageTypeVO, err := asset.NewImageType(string(graphQLImage.Type))
		if err != nil {
			return nil, err
		}
		imageType = imageTypeVO
	}

	var storageLocation *asset.S3ObjectValue
	if graphQLImage.StorageLocation != nil {
		storageLocationVO, err := asset.NewS3ObjectValue(
			graphQLImage.StorageLocation.Bucket,
			graphQLImage.StorageLocation.Key,
			graphQLImage.StorageLocation.URL,
		)
		if err != nil {
			return nil, err
		}
		storageLocation = storageLocationVO
	}

	metadata := convertStringSliceToString(graphQLImage.Metadata)

	return asset.NewImage(
		*imageID,
		*fileName,
		graphQLImage.URL,
		imageType,
		storageLocation,
		graphQLImage.Width,
		graphQLImage.Height,
		graphQLImage.Size,
		graphQLImage.ContentType,
		nil, // streamInfo - not available in GraphQL
		metadata,
		graphQLImage.CreatedAt,
		graphQLImage.UpdatedAt,
	), nil
}

type GraphQLAsset struct {
	ID          string              `json:"id"`
	Slug        string              `json:"slug"`
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	Type        string              `json:"type"`
	Genre       *string             `json:"genre"`
	Genres      []string            `json:"genres"`
	Tags        []string            `json:"tags"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Metadata    *string             `json:"metadata"`
	OwnerID     *string             `json:"ownerId"`
	Videos      []GraphQLVideo      `json:"videos"`
	Images      []GraphQLImage      `json:"images"`
	PublishRule *GraphQLPublishRule `json:"publishRule"`
}

type GraphQLVideo struct {
	ID                 string                  `json:"id"`
	Label              string                  `json:"label"`
	Type               VideoType               `json:"type"`
	Format             VideoFormat             `json:"format"`
	StorageLocation    GraphQLS3Object         `json:"storageLocation"`
	Width              *int                    `json:"width"`
	Height             *int                    `json:"height"`
	Duration           *float64                `json:"duration"`
	Bitrate            *int                    `json:"bitrate"`
	Codec              *string                 `json:"codec"`
	Size               *int                    `json:"size"`
	ContentType        *string                 `json:"contentType"`
	StreamInfo         *GraphQLStreamInfo      `json:"streamInfo"`
	Metadata           []string                `json:"metadata"`
	Status             VideoStatus             `json:"status"`
	Thumbnail          *GraphQLImage           `json:"thumbnail"`
	CreatedAt          time.Time               `json:"createdAt"`
	UpdatedAt          time.Time               `json:"updatedAt"`
	Quality            *string                 `json:"quality"`
	IsReady            bool                    `json:"isReady"`
	IsProcessing       bool                    `json:"isProcessing"`
	IsFailed           bool                    `json:"isFailed"`
	SegmentCount       *int                    `json:"segmentCount"`
	VideoCodec         *string                 `json:"videoCodec"`
	AudioCodec         *string                 `json:"audioCodec"`
	AvgSegmentDuration *float64                `json:"avgSegmentDuration"`
	Segments           []string                `json:"segments"`
	FrameRate          *string                 `json:"frameRate"`
	AudioChannels      *int                    `json:"audioChannels"`
	AudioSampleRate    *int                    `json:"audioSampleRate"`
	TranscodingInfo    *GraphQLTranscodingInfo `json:"transcodingInfo"`
}

type GraphQLImage struct {
	ID              string           `json:"id"`
	FileName        string           `json:"fileName"`
	URL             string           `json:"url"`
	Type            ImageType        `json:"type"`
	StorageLocation *GraphQLS3Object `json:"storageLocation"`
	Width           *int             `json:"width"`
	Height          *int             `json:"height"`
	Size            *int             `json:"size"`
	ContentType     *string          `json:"contentType"`
	Metadata        []string         `json:"metadata"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type GraphQLS3Object struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type GraphQLStreamInfo struct {
	DownloadURL *string `json:"downloadUrl"`
	CDNPrefix   *string `json:"cdnPrefix"`
	URL         *string `json:"url"`
}

type GraphQLPublishRule struct {
	PublishAt   *time.Time `json:"publishAt"`
	UnpublishAt *time.Time `json:"unpublishAt"`
	Regions     []string   `json:"regions"`
	AgeRating   *string    `json:"ageRating"`
}

type GraphQLTranscodingInfo struct {
	JobID       *string    `json:"jobId"`
	Progress    *float64   `json:"progress"`
	OutputURL   *string    `json:"outputUrl"`
	Error       *string    `json:"error"`
	CompletedAt *time.Time `json:"completedAt"`
}

type VideoType string
type VideoFormat string
type VideoStatus string
type ImageType string

func convertStringSliceToString(slice []string) *string {
	if len(slice) == 0 {
		return nil
	}
	// TODO: Consider JSON marshaling for complex metadata
	result := strings.Join(slice, ",")
	return &result
}
