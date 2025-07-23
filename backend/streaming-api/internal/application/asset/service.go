package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	assetdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	bucketdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type CacheService interface {
	GetAsset(ctx context.Context, slug string) (*assetdomain.Asset, error)
	SetAsset(ctx context.Context, asset *assetdomain.Asset) error
	GetAssets(ctx context.Context) ([]*assetdomain.Asset, error)
	SetAssets(ctx context.Context, assets []*assetdomain.Asset) error
}

type ApplicationService struct {
	repo              assetdomain.Repository
	cacheService      CacheService
	publishingService *assetdomain.AssetPublishingService
	streamingService  *assetdomain.AssetStreamingService
	searchService     *assetdomain.AssetSearchService
	logger            *logger.Logger
	circuitBreaker    *errors.CircuitBreaker
}

func NewApplicationService(
	repo assetdomain.Repository,
	cacheService CacheService,
	circuitBreaker *errors.CircuitBreaker,
) *ApplicationService {
	return &ApplicationService{
		repo:              repo,
		cacheService:      cacheService,
		publishingService: assetdomain.NewAssetPublishingService(),
		streamingService:  assetdomain.NewAssetStreamingService(),
		searchService:     assetdomain.NewAssetSearchService(),
		logger:            logger.Get().WithService("asset-application"),
		circuitBreaker:    circuitBreaker,
	}
}

func (s *ApplicationService) GetAsset(ctx context.Context, slug assetdomain.Slug) (*assetdomain.Asset, error) {
	asset, err := s.cacheService.GetAsset(ctx, slug.Value())
	if err != nil {
		s.logger.WithError(err).Error("Failed to get asset from cache", "slug", slug.Value())
		return nil, errors.NewTransientError("cache error", err)
	}

	if asset != nil {
		s.logger.Debug("Asset found in cache", "slug", slug.Value())
		return asset, nil
	}

	s.logger.Debug("Asset not found in cache, fetching from repository", "slug", slug.Value())

	asset, err = s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_asset",
			"slug":      slug.Value(),
		})
	}

	if asset != nil {
		if err := s.cacheService.SetAsset(ctx, asset); err != nil {
			s.logger.WithError(err).Warn("Failed to cache asset", "slug", slug.Value())
		}
	}

	return asset, nil
}

func (s *ApplicationService) GetAssets(ctx context.Context) ([]*assetdomain.Asset, error) {
	assets, err := s.cacheService.GetAssets(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get assets from cache")
		return nil, errors.NewTransientError("cache error", err)
	}

	if assets != nil {
		s.logger.Debug("Assets found in cache", "count", len(assets))
		return assets, nil
	}

	s.logger.Debug("Assets not found in cache, fetching from repository")

	assets, err = s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_assets",
		})
	}

	if assets != nil {
		if err := s.cacheService.SetAssets(ctx, assets); err != nil {
			s.logger.WithError(err).Warn("Failed to cache assets")
		}
	}

	return assets, nil
}

func (s *ApplicationService) GetPublicAssets(ctx context.Context) ([]*assetdomain.Asset, error) {
	assets, err := s.repo.GetPublic(ctx)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_public_assets",
		})
	}

	return assets, nil
}

func (s *ApplicationService) GetAssetsByType(ctx context.Context, assetType assetdomain.AssetType) ([]*assetdomain.Asset, error) {
	assets, err := s.repo.GetByType(ctx, assetType)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_assets_by_type",
			"type":      assetType.Value(),
		})
	}

	return assets, nil
}

func (s *ApplicationService) GetAssetsByGenre(ctx context.Context, genre assetdomain.Genre) ([]*assetdomain.Asset, error) {
	assets, err := s.repo.GetByGenre(ctx, genre)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_assets_by_genre",
			"genre":     genre.Value(),
		})
	}

	return assets, nil
}

func (s *ApplicationService) GetAssetsInBucket(ctx context.Context, bucketKey bucketdomain.BucketKey) ([]*assetdomain.Asset, error) {
	assets, err := s.GetAssets(ctx)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *ApplicationService) SearchAssets(ctx context.Context, query string, filters *assetdomain.SearchFilters) ([]*assetdomain.Asset, error) {
	assets, err := s.searchService.SearchAssets(query, filters)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "search_assets",
			"query":     query,
		})
	}

	return assets, nil
}

func (s *ApplicationService) GetStreamingInfo(ctx context.Context, slug assetdomain.Slug, userID string, region string, userAge int) (*assetdomain.StreamingInfo, error) {
	asset, err := s.GetAsset(ctx, slug)
	if err != nil {
		return nil, err
	}

	streamingInfo, err := s.streamingService.GetStreamingInfo(asset, userID, region, userAge)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "get_streaming_info",
			"slug":      slug.Value(),
			"userID":    userID,
			"region":    region,
		})
	}

	return streamingInfo, nil
}

func (s *ApplicationService) GetRecommendedAssets(ctx context.Context, slug assetdomain.Slug, limit int) ([]*assetdomain.Asset, error) {
	asset, err := s.GetAsset(ctx, slug)
	if err != nil {
		return nil, err
	}

	recommendations := s.streamingService.GetRecommendedAssets(asset, limit)
	return recommendations, nil
}

func (s *ApplicationService) GetPublishStatus(ctx context.Context, slug assetdomain.Slug) (constants.PublishStatus, error) {
	asset, err := s.GetAsset(ctx, slug)
	if err != nil {
		return constants.PublishStatusInvalid, err
	}

	status := s.publishingService.GetPublishStatus(asset)
	return status, nil
}
