package graphql

import (
	"context"
	"strings"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	resilience "github.com/serdarburakguneri/hobby-streamer/backend/pkg/resilience"
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
	bucketvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/graphql/queries"
)

type BucketRepository struct {
	client         *Client
	logger         *logger.Logger
	circuitBreaker *resilience.CircuitBreaker
}

func NewBucketRepository(client *Client, circuitBreaker *resilience.CircuitBreaker) bucket.Repository {
	return &BucketRepository{
		client:         client,
		logger:         logger.Get().WithService("bucket-graphql-repository"),
		circuitBreaker: circuitBreaker,
	}
}

func (r *BucketRepository) GetByID(ctx context.Context, id bucketvalueobjects.BucketID) (*bucketentity.Bucket, error) {
	buckets, err := r.GetAll(ctx, 10, nil)
	if err != nil {
		return nil, err
	}

	for _, b := range buckets {
		if b.ID().Equals(id) {
			return b, nil
		}
	}

	return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
}

func (r *BucketRepository) GetByKey(ctx context.Context, key bucketvalueobjects.BucketKey) (*bucketentity.Bucket, error) {
	query := queries.GetBucketQuery

	variables := map[string]interface{}{"key": key.Value()}

	var response struct {
		BucketByKey *GraphQLBucketWithAssets `json:"bucketByKey"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, variables, &response)
	})

	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{"operation": "get_bucket", "key": key.Value()})
	}

	if response.BucketByKey == nil {
		return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
	}

	return r.convertGraphQLBucketWithAssetsToDomain(response.BucketByKey)
}

func (r *BucketRepository) GetAll(ctx context.Context, limit int, nextKey *string) ([]*bucketentity.Bucket, error) {
	query := queries.GetBucketsQuery

	variables := map[string]interface{}{"limit": limit, "nextKey": nextKey}

	var response struct {
		Buckets struct {
			Items []*GraphQLBucketWithAssets `json:"items"`
		} `json:"buckets"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, variables, &response)
	})

	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{"operation": "get_buckets"})
	}

	var buckets []*bucketentity.Bucket
	for _, graphQLBucket := range response.Buckets.Items {
		bkt, err := r.convertGraphQLBucketWithAssetsToDomain(graphQLBucket)
		if err != nil {
			r.logger.WithError(err).Error("Failed to convert GraphQL bucket to domain", "bucket_id", graphQLBucket.ID)
			continue
		}
		buckets = append(buckets, bkt)
	}

	return buckets, nil
}

func (r *BucketRepository) GetByType(ctx context.Context, bucketType bucketvalueobjects.BucketType) ([]*bucketentity.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var filteredBuckets []*bucketentity.Bucket
	for _, b := range buckets {
		if b.Type().Equals(bucketType) {
			filteredBuckets = append(filteredBuckets, b)
		}
	}

	return filteredBuckets, nil
}

func (r *BucketRepository) GetActive(ctx context.Context) ([]*bucketentity.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var activeBuckets []*bucketentity.Bucket
	for _, b := range buckets {
		if b.IsActive() {
			activeBuckets = append(activeBuckets, b)
		}
	}

	return activeBuckets, nil
}

func (r *BucketRepository) GetByAssetType(ctx context.Context, assetType assetvalueobjects.AssetType) ([]*bucketentity.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var filteredBuckets []*bucketentity.Bucket
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

func (r *BucketRepository) GetByGenre(ctx context.Context, genre assetvalueobjects.Genre) ([]*bucketentity.Bucket, error) {
	buckets, err := r.GetAll(context.Background(), 10, nil)
	if err != nil {
		return nil, err
	}

	var filteredBuckets []*bucketentity.Bucket
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

func (r *BucketRepository) GetAssets(ctx context.Context, bucket *bucketentity.Bucket) ([]*assetentity.Asset, error) {
	return bucket.Assets(), nil
}

func (r *BucketRepository) GetRecommendations(ctx context.Context, bucket *bucketentity.Bucket, limit int) ([]*assetentity.Asset, error) {
	return nil, pkgerrors.NewNotFoundError("recommendations not implemented", nil)
}

func (r *BucketRepository) convertGraphQLBucketWithAssetsToDomain(graphQLBucket *GraphQLBucketWithAssets) (*bucketentity.Bucket, error) {
	bucketID, err := bucketvalueobjects.NewBucketID(graphQLBucket.ID)
	if err != nil {
		return nil, err
	}

	bucketKey, err := bucketvalueobjects.NewBucketKey(graphQLBucket.Key)
	if err != nil {
		return nil, err
	}

	bucketName, err := bucketvalueobjects.NewBucketName(graphQLBucket.Name)
	if err != nil {
		return nil, err
	}

	var description *bucketvalueobjects.BucketDescription
	if graphQLBucket.Description != nil {
		descVO, err := bucketvalueobjects.NewBucketDescription(*graphQLBucket.Description)
		if err != nil {
			return nil, err
		}
		description = descVO
	}

	bucketType, err := bucketvalueobjects.NewBucketType(strings.ToLower(graphQLBucket.Type))
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{"raw_type": graphQLBucket.Type})
	}

	var status *bucketvalueobjects.BucketStatus
	if graphQLBucket.Status != nil {
		statusVO, err := bucketvalueobjects.NewBucketStatus(*graphQLBucket.Status)
		if err != nil {
			return nil, err
		}
		status = statusVO
	}

	createdAt := bucketvalueobjects.NewCreatedAt(graphQLBucket.CreatedAt)
	updatedAt := bucketvalueobjects.NewUpdatedAt(graphQLBucket.UpdatedAt)

	assetIDs := make([]string, len(graphQLBucket.Assets))
	for i, asset := range graphQLBucket.Assets {
		assetIDs[i] = asset.ID
	}

	assetIDsVO, err := bucketvalueobjects.NewAssetIDs(assetIDs)
	if err != nil {
		return nil, err
	}

	assets, err := ConvertGraphQLAssetsToDomain(graphQLBucket.Assets)
	if err != nil {
		return nil, err
	}

	return bucketentity.NewBucket(
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
