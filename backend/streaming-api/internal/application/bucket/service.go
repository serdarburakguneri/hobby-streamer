package bucket

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	assetdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	bucketdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type CacheService interface {
	GetBucket(ctx context.Context, key string) (*bucketdomain.Bucket, error)
	SetBucket(ctx context.Context, bucket *bucketdomain.Bucket) error
	GetBuckets(ctx context.Context) ([]*bucketdomain.Bucket, error)
	SetBuckets(ctx context.Context, buckets []*bucketdomain.Bucket) error
}

type ApplicationService struct {
	repo                bucketdomain.Repository
	cacheService        CacheService
	organizationService *bucketdomain.BucketOrganizationService
	discoveryService    *bucketdomain.BucketDiscoveryService
	logger              *logger.Logger
	circuitBreaker      *errors.CircuitBreaker
}

func NewApplicationService(
	repo bucketdomain.Repository,
	cacheService CacheService,
	circuitBreaker *errors.CircuitBreaker,
) *ApplicationService {
	return &ApplicationService{
		repo:                repo,
		cacheService:        cacheService,
		organizationService: bucketdomain.NewBucketOrganizationService(),
		discoveryService:    bucketdomain.NewBucketDiscoveryService(),
		logger:              logger.Get().WithService("bucket-application"),
		circuitBreaker:      circuitBreaker,
	}
}

func (s *ApplicationService) GetBucket(ctx context.Context, key bucketdomain.BucketKey) (*bucketdomain.Bucket, error) {
	bucket, err := s.cacheService.GetBucket(ctx, key.Value())
	if err != nil {
		s.logger.WithError(err).Error("Failed to get bucket from cache", "key", key.Value())
		return nil, errors.NewTransientError("cache error", err)
	}

	if bucket != nil {
		s.logger.Debug("Bucket found in cache", "key", key.Value())
		return bucket, nil
	}

	s.logger.Debug("Bucket not found in cache, fetching from repository", "key", key.Value())

	bucket, err = s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_bucket",
			"key":       key.Value(),
		})
	}

	if bucket != nil {
		if err := s.cacheService.SetBucket(ctx, bucket); err != nil {
			s.logger.WithError(err).Warn("Failed to cache bucket", "key", key.Value())
		}
	}

	return bucket, nil
}

func (s *ApplicationService) GetBuckets(ctx context.Context) ([]*bucketdomain.Bucket, error) {
	buckets, err := s.cacheService.GetBuckets(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get buckets from cache")
		return nil, errors.NewTransientError("cache error", err)
	}

	if buckets != nil {
		s.logger.Debug("Buckets found in cache", "count", len(buckets))
		return buckets, nil
	}

	s.logger.Debug("Buckets not found in cache, fetching from repository")

	buckets, err = s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_buckets",
		})
	}

	if buckets != nil {
		if err := s.cacheService.SetBuckets(ctx, buckets); err != nil {
			s.logger.WithError(err).Warn("Failed to cache buckets")
		}
	}

	return buckets, nil
}

func (s *ApplicationService) GetBucketsByType(ctx context.Context, bucketType string) ([]*bucketdomain.Bucket, error) {
	bucketTypeVO, err := bucketdomain.NewBucketType(bucketType)
	if err != nil {
		return nil, errors.NewValidationError("invalid bucket type", err)
	}

	buckets, err := s.repo.GetByType(ctx, *bucketTypeVO)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_buckets_by_type",
			"type":      bucketType,
		})
	}

	return buckets, nil
}

func (s *ApplicationService) GetActiveBuckets(ctx context.Context) ([]*bucketdomain.Bucket, error) {
	buckets, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_active_buckets",
		})
	}

	return buckets, nil
}

func (s *ApplicationService) GetBucketStats(ctx context.Context, key bucketdomain.BucketKey) (*bucketdomain.BucketStats, error) {
	bucket, err := s.GetBucket(ctx, key)
	if err != nil {
		return nil, err
	}

	stats := s.organizationService.GetBucketStats(bucket)
	return stats, nil
}

func (s *ApplicationService) GetRecommendedAssets(ctx context.Context, key bucketdomain.BucketKey, limit int) ([]*assetdomain.Asset, error) {
	bucket, err := s.GetBucket(ctx, key)
	if err != nil {
		return nil, err
	}

	recommendations := s.organizationService.GetBucketRecommendations(bucket, limit)
	return recommendations, nil
}

func (s *ApplicationService) SearchBuckets(ctx context.Context, query string, filters *bucketdomain.BucketSearchFilters) ([]*bucketdomain.Bucket, error) {
	buckets, err := s.discoveryService.SearchBuckets(query, filters)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "search_buckets",
			"query":     query,
		})
	}

	return buckets, nil
}
