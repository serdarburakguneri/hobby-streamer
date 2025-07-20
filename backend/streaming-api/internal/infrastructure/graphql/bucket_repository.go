package graphql

import (
	"context"
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
	buckets, err := r.GetAll(ctx)
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
					ownerId
					createdAt
					updatedAt
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

func (r *BucketRepository) GetAll(ctx context.Context) ([]*bucket.Bucket, error) {
	query := `
		query GetBuckets {
			buckets {
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
						ownerId
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	var response struct {
		Buckets struct {
			Items []GraphQLBucketWithAssets `json:"items"`
		} `json:"buckets"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, map[string]interface{}{}, &response)
	})

	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_buckets",
		})
	}

	buckets := make([]*bucket.Bucket, len(response.Buckets.Items))
	for i, graphQLBucket := range response.Buckets.Items {
		domainBucket, err := r.convertGraphQLBucketWithAssetsToDomain(&graphQLBucket)
		if err != nil {
			r.logger.WithError(err).Error("Failed to convert GraphQL bucket to domain", "bucket_id", graphQLBucket.ID)
			continue
		}
		buckets[i] = domainBucket
	}

	return buckets, nil
}

func (r *BucketRepository) GetByType(ctx context.Context, bucketType bucket.BucketType) ([]*bucket.Bucket, error) {
	buckets, err := r.GetAll(ctx)
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
	buckets, err := r.GetAll(ctx)
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
	buckets, err := r.GetAll(ctx)
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
	buckets, err := r.GetAll(ctx)
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
		[]asset.Video{}, // Videos not loaded in bucket context
		[]asset.Image{}, // Images not loaded in bucket context
		nil,             // Publish rule not loaded in bucket context
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
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Type        string    `json:"type"`
	Genre       *string   `json:"genre"`
	Genres      []string  `json:"genres"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	Metadata    *string   `json:"metadata"`
	OwnerID     *string   `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
