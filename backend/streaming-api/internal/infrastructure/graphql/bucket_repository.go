package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type BucketRepository struct {
	client         *Client
	logger         *logger.Logger
	circuitBreaker *errors.CircuitBreaker
}

func NewBucketRepository(client *Client, circuitBreaker *errors.CircuitBreaker) bucket.Repository {
	return &BucketRepository{
		client:         client,
		logger:         logger.Get().WithService("bucket-graphql-repository"),
		circuitBreaker: circuitBreaker,
	}
}

func (r *BucketRepository) GetByID(ctx context.Context, id bucket.BucketID) (*bucket.Bucket, error) {
	buckets, err := r.GetAll(ctx, 10, nil)
	if err != nil {
		return nil, err
	}

	for _, b := range buckets {
		if b.ID().Equals(id) {
			return b, nil
		}
	}

	return nil, errors.NewNotFoundError("bucket not found", nil)
}

func (r *BucketRepository) GetByKey(ctx context.Context, key bucket.BucketKey) (*bucket.Bucket, error) {
	query := `
		query GetBucket($key: String!) {
			bucketByKey(key: $key) {
				id
				key
				name
				description
				type
				status
				createdAt
				updatedAt
				assets {
					id
					slug
					title
					description
					type
					genre
					genres
					tags
					status
					metadata
					createdAt
					updatedAt
					ownerId
					videos {
						id
						type
						format
						storageLocation { bucket key url }
						width
						height
						duration
						bitrate
						codec
						size
						contentType
						streamInfo { downloadUrl cdnPrefix url }
						metadata
						status
						thumbnail {
							id
							fileName
							url
							type
							storageLocation { bucket key url }
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
						transcodingInfo { jobId progress outputUrl error completedAt }
					}
					images {
						id
						fileName
						url
						type
						storageLocation { bucket key url }
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

	variables := map[string]interface{}{
		"key": key.Value(),
	}

	var response struct {
		BucketByKey *GraphQLBucketWithAssets `json:"bucketByKey"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, variables, &response)
	})

	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_bucket",
			"key":       key.Value(),
		})
	}

	if response.BucketByKey == nil {
		return nil, errors.NewNotFoundError("bucket not found", nil)
	}

	return r.convertGraphQLBucketWithAssetsToDomain(response.BucketByKey)
}

func (r *BucketRepository) GetAll(ctx context.Context, limit int, nextKey *string) ([]*bucket.Bucket, error) {
	query := `
		query GetBuckets($limit: Int, $nextKey: String) {
			buckets(limit: $limit, nextKey: $nextKey) {
				items {
					id
					key
					name
					description
					type
					status
					createdAt
					updatedAt
					assets {
						id
						slug
						title
						description
						type
						genre
						genres
						tags
						status
						metadata
						createdAt
						updatedAt
						ownerId
						videos {
							id
							type
							format
							storageLocation { bucket key url }
							width
							height
							duration
							bitrate
							codec
							size
							contentType
							streamInfo { downloadUrl cdnPrefix url }
							metadata
							status
							thumbnail {
								id
								fileName
								url
								type
								storageLocation { bucket key url }
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
							transcodingInfo { jobId progress outputUrl error completedAt }
						}
						images {
							id
							fileName
							url
							type
							storageLocation { bucket key url }
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
				nextKey
				hasMore
			}
		}
	`

	variables := map[string]interface{}{
		"limit": limit,
	}
	if nextKey != nil {
		variables["nextKey"] = *nextKey
	}

	var response struct {
		Buckets struct {
			Items   []GraphQLBucketWithAssets `json:"items"`
			NextKey *string                   `json:"nextKey"`
			HasMore bool                      `json:"hasMore"`
		} `json:"buckets"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, variables, &response)
	})

	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_buckets",
		})
	}

	var buckets []*bucket.Bucket
	for _, graphQLBucket := range response.Buckets.Items {
		domainBucket, err := r.convertGraphQLBucketWithAssetsToDomain(&graphQLBucket)
		if err != nil {
			r.logger.WithError(err).Error("Failed to convert GraphQL bucket to domain", "bucket_id", graphQLBucket.ID, "raw", fmt.Sprintf("%+v", graphQLBucket))
			continue
		}
		if domainBucket == nil {
			r.logger.Error("convertGraphQLBucketWithAssetsToDomain returned nil", "bucket_id", graphQLBucket.ID, "raw", fmt.Sprintf("%+v", graphQLBucket))
			continue
		}
		buckets = append(buckets, domainBucket)
	}

	return buckets, nil
}

func (r *BucketRepository) GetByType(ctx context.Context, bucketType bucket.BucketType) ([]*bucket.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var filteredBuckets []*bucket.Bucket
	for _, b := range buckets {
		if b.Type().Equals(bucketType) {
			filteredBuckets = append(filteredBuckets, b)
		}
	}

	return filteredBuckets, nil
}

func (r *BucketRepository) GetActive(ctx context.Context) ([]*bucket.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var activeBuckets []*bucket.Bucket
	for _, b := range buckets {
		if b.IsActive() {
			activeBuckets = append(activeBuckets, b)
		}
	}

	return activeBuckets, nil
}

func (r *BucketRepository) GetByAssetType(ctx context.Context, assetType asset.AssetType) ([]*bucket.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var filteredBuckets []*bucket.Bucket
	for _, b := range buckets {
		assets, err := r.GetAssets(ctx, b)
		if err != nil {
			continue
		}
		for _, a := range assets {
			if a.Type().Equals(assetType) {
				filteredBuckets = append(filteredBuckets, b)
				break
			}
		}
	}

	return filteredBuckets, nil
}

func (r *BucketRepository) GetByGenre(ctx context.Context, genre asset.Genre) ([]*bucket.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var filteredBuckets []*bucket.Bucket
	for _, b := range buckets {
		assets, err := r.GetAssets(ctx, b)
		if err != nil {
			continue
		}
		for _, a := range assets {
			if a.Genre() != nil && a.Genre().Equals(genre) {
				filteredBuckets = append(filteredBuckets, b)
				break
			}
		}
	}

	return filteredBuckets, nil
}

func (r *BucketRepository) GetAssets(ctx context.Context, bucket *bucket.Bucket) ([]*asset.Asset, error) {
	return bucket.Assets(), nil
}

func (r *BucketRepository) Search(ctx context.Context, query string, filters *bucket.BucketSearchFilters) ([]*bucket.Bucket, error) {
	return nil, errors.NewNotFoundError("search not implemented", nil)
}

func (r *BucketRepository) GetStats(ctx context.Context, bucket *bucket.Bucket) (*bucket.BucketStats, error) {
	return nil, errors.NewNotFoundError("stats not implemented", nil)
}

func (r *BucketRepository) GetRecommendations(ctx context.Context, bucket *bucket.Bucket, limit int) ([]*asset.Asset, error) {
	return nil, errors.NewNotFoundError("recommendations not implemented", nil)
}

func (r *BucketRepository) convertGraphQLBucketWithAssetsToDomain(graphQLBucket *GraphQLBucketWithAssets) (*bucket.Bucket, error) {
	bucketID, err := bucket.NewBucketID(graphQLBucket.ID)
	if err != nil {
		return nil, err
	}

	bucketKey, err := bucket.NewBucketKey(graphQLBucket.Key)
	if err != nil {
		return nil, err
	}

	bucketName, err := bucket.NewBucketName(graphQLBucket.Name)
	if err != nil {
		return nil, err
	}

	var description *bucket.BucketDescription
	if graphQLBucket.Description != nil {
		descVO, err := bucket.NewBucketDescription(*graphQLBucket.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	bucketType, err := bucket.NewBucketType(graphQLBucket.Type)
	if err != nil {
		return nil, err
	}

	var status *bucket.BucketStatus
	if graphQLBucket.Status != nil {
		statusVO, err := bucket.NewBucketStatus(*graphQLBucket.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := bucket.NewCreatedAt(graphQLBucket.CreatedAt)
	updatedAt := bucket.NewUpdatedAt(graphQLBucket.UpdatedAt)

	assetIDs := make([]string, len(graphQLBucket.Assets))
	for i, asset := range graphQLBucket.Assets {
		assetIDs[i] = asset.ID
	}

	assetIDsVO, err := bucket.NewAssetIDs(assetIDs)
	if err != nil {
		return nil, err
	}

	assets, err := r.convertGraphQLAssetsToDomain(graphQLBucket.Assets)
	if err != nil {
		return nil, err
	}

	return bucket.NewBucket(
		*bucketID,
		*bucketKey,
		*bucketName,
		description,
		*bucketType,
		status,
		assetIDsVO,
		*createdAt,
		*updatedAt,
		assets,
	), nil
}

func (r *BucketRepository) convertGraphQLAssetsToDomain(graphQLAssets []GraphQLBucketAsset) ([]*asset.Asset, error) {
	assets := make([]*asset.Asset, len(graphQLAssets))
	for i, graphQLAsset := range graphQLAssets {
		domainAsset, err := r.convertGraphQLAssetToDomain(graphQLAsset)
		if err != nil {
			return nil, err
		}
		assets[i] = domainAsset
	}
	return assets, nil
}

func (r *BucketRepository) convertGraphQLAssetToDomain(graphQLAsset GraphQLBucketAsset) (*asset.Asset, error) {
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

func (r *BucketRepository) convertGraphQLVideosToDomain(graphQLVideos []GraphQLVideo) ([]asset.Video, error) {
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

func (r *BucketRepository) convertGraphQLImagesToDomain(graphQLImages []GraphQLImage) ([]asset.Image, error) {
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

func (r *BucketRepository) convertGraphQLVideoToDomain(graphQLVideo GraphQLVideo) (*asset.Video, error) {
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

func (r *BucketRepository) convertGraphQLImageToDomain(graphQLImage *GraphQLImage) (*asset.Image, error) {
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

type GraphQLBucket struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Type        string    `json:"type"`
	Status      *string   `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GraphQLBucketWithAssets struct {
	ID          string               `json:"id"`
	Key         string               `json:"key"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Type        string               `json:"type"`
	Status      *string              `json:"status"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Assets      []GraphQLBucketAsset `json:"assets"`
}

type GraphQLBucketAsset struct {
	ID          string              `json:"id"`
	Slug        string              `json:"slug"`
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	Type        string              `json:"type"`
	Genre       *string             `json:"genre"`
	Genres      []string            `json:"genres"`
	Tags        []string            `json:"tags"`
	Status      string              `json:"status"`
	Metadata    *string             `json:"metadata"`
	OwnerID     *string             `json:"ownerId"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Videos      []GraphQLVideo      `json:"videos"`
	Images      []GraphQLImage      `json:"images"`
	PublishRule *GraphQLPublishRule `json:"publishRule"`
}
