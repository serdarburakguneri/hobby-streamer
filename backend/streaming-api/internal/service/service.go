package service

import (
	"context"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/cache"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/model"
)

type ServiceInterface interface {
	GetBuckets(ctx context.Context) ([]model.Bucket, error)
	GetBucket(ctx context.Context, key string) (*model.Bucket, error)
	GetAssets(ctx context.Context) ([]model.Asset, error)
	GetAsset(ctx context.Context, slug string) (*model.Asset, error)
	GetAssetsInBucket(ctx context.Context, bucketKey string) ([]model.Asset, error)
}

type Service struct {
	cacheService    cache.CacheService
	logger          *logger.Logger
	assetManagerURL string
	serviceClient   auth.ServiceClientInterface
	circuitBreaker  *errors.CircuitBreaker
	graphQLClient   *GraphQLClient
}

func NewService(cacheService cache.CacheService, assetManagerURL, keycloakURL, realm, clientID, clientSecret string, circuitBreakerConfig errors.CircuitBreakerConfig) *Service {
	log := logger.Get().WithService("streaming-service")
	log.Info("Initializing streaming service", "asset_manager_url", assetManagerURL, "keycloak_url", keycloakURL, "realm", realm, "client_id", clientID)

	serviceClient := auth.NewServiceClient(keycloakURL, realm, clientID, clientSecret)
	log.Info("Service client created successfully")

	circuitBreaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
		Name:      "asset-manager",
		Threshold: circuitBreakerConfig.Threshold,
		Timeout:   circuitBreakerConfig.Timeout,
		OnStateChange: func(name string, from, to errors.CircuitState) {
			log.Info("Circuit breaker state changed", "name", name, "from", from, "to", to)
		},
	})

	return &Service{
		cacheService:    cacheService,
		logger:          log,
		assetManagerURL: assetManagerURL,
		serviceClient:   serviceClient,
		circuitBreaker:  circuitBreaker,
		graphQLClient:   NewGraphQLClient(serviceClient),
	}
}

func (s *Service) GetBucket(ctx context.Context, key string) (*model.Bucket, error) {
	bucket, err := s.cacheService.GetBucket(ctx, key)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get bucket from cache", "key", key)
		return nil, errors.NewTransientError("cache error", err)
	}

	if bucket != nil {
		s.logger.Debug("Bucket found in cache", "key", key)
		return bucket, nil
	}

	s.logger.Debug("Bucket not found in cache, fetching from asset-manager", "key", key)

	bucket, err = s.fetchBucketFromAssetManager(ctx, key)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "fetch_bucket",
			"key":       key,
		})
	}

	if bucket != nil {
		if err := s.cacheService.SetBucket(ctx, bucket); err != nil {
			s.logger.WithError(err).Warn("Failed to cache bucket", "key", key)
		}
	}

	return bucket, nil
}

func (s *Service) GetBuckets(ctx context.Context) ([]model.Bucket, error) {
	buckets, err := s.cacheService.GetBuckets(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get buckets from cache")
		return nil, errors.NewTransientError("cache error", err)
	}

	if buckets != nil {
		s.logger.Debug("Buckets found in cache", "count", len(buckets))
		return buckets, nil
	}

	s.logger.Debug("Buckets not found in cache, fetching from asset-manager")

	buckets, err = s.fetchBucketsFromAssetManager(ctx)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "fetch_buckets",
		})
	}

	if buckets != nil {
		if err := s.cacheService.SetBuckets(ctx, buckets); err != nil {
			s.logger.WithError(err).Warn("Failed to cache buckets")
		}
	}

	return buckets, nil
}

func (s *Service) GetAsset(ctx context.Context, slug string) (*model.Asset, error) {
	asset, err := s.cacheService.GetAsset(ctx, slug)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get asset from cache", "slug", slug)
		return nil, errors.NewTransientError("cache error", err)
	}

	if asset != nil {
		s.logger.Debug("Asset found in cache", "slug", slug)
		return asset, nil
	}

	s.logger.Debug("Asset not found in cache, fetching from asset-manager", "slug", slug)

	asset, err = s.fetchAssetFromAssetManager(ctx, slug)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "fetch_asset",
			"slug":      slug,
		})
	}

	if asset != nil {
		if err := s.cacheService.SetAsset(ctx, asset); err != nil {
			s.logger.WithError(err).Warn("Failed to cache asset", "slug", slug)
		}
	}

	return asset, nil
}

func (s *Service) GetAssets(ctx context.Context) ([]model.Asset, error) {
	assets, err := s.cacheService.GetAssets(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get assets from cache")
		return nil, errors.NewTransientError("cache error", err)
	}

	if assets != nil {
		s.logger.Debug("Assets found in cache", "count", len(assets))
		return assets, nil
	}

	s.logger.Debug("Assets not found in cache, fetching from asset-manager")

	assets, err = s.fetchAssetsFromAssetManager(ctx)
	if err != nil {
		return nil, errors.WithContext(err, map[string]interface{}{
			"operation": "fetch_assets",
		})
	}

	if assets != nil {
		if err := s.cacheService.SetAssets(ctx, assets); err != nil {
			s.logger.WithError(err).Warn("Failed to cache assets")
		}
	}

	return assets, nil
}

func (s *Service) GetAssetsInBucket(ctx context.Context, bucketKey string) ([]model.Asset, error) {
	bucket, err := s.GetBucket(ctx, bucketKey)
	if err != nil {
		return nil, errors.NewNotFoundError("failed to get bucket", err)
	}

	if bucket == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("bucket not found: %s", bucketKey), nil)
	}

	if len(bucket.AssetIDs) == 0 {
		return []model.Asset{}, nil
	}

	allAssets, err := s.GetAssets(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to fetch all assets for bucket filtering", "bucket_key", bucketKey)
		return nil, errors.NewExternalError("failed to fetch assets", err)
	}

	var bucketAssets []model.Asset
	assetIDSet := make(map[string]bool)
	for _, id := range bucket.AssetIDs {
		assetIDSet[id] = true
	}

	for _, asset := range allAssets {
		if assetIDSet[asset.ID] {
			bucketAssets = append(bucketAssets, asset)
		}
	}

	s.logger.Debug("Filtered assets for bucket", "bucket_key", bucketKey, "total_assets", len(allAssets), "bucket_asset_ids", len(bucket.AssetIDs), "filtered_assets", len(bucketAssets))

	return bucketAssets, nil
}

func (s *Service) fetchBucketFromAssetManager(ctx context.Context, key string) (*model.Bucket, error) {
	query := fmt.Sprintf(`
		query {
			bucketByKey(key: "%s") {
				id
				key
				name
				description
				type
				status
				assetIds
				createdAt
				updatedAt
			}
		}
	`, key)

	var response struct {
		Data struct {
			BucketByKey *model.Bucket `json:"bucketByKey"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	err := s.graphQLClient.ExecuteQueryWithCircuitBreaker(ctx, s.circuitBreaker, s.assetManagerURL+"/graphql", query, &response)

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch bucket from asset-manager", err)
	}

	if err := s.graphQLClient.HandleGraphQLErrors(&response); err != nil {
		return nil, err
	}

	return response.Data.BucketByKey, nil
}

func (s *Service) fetchBucketsFromAssetManager(ctx context.Context) ([]model.Bucket, error) {
	query := `
		query {
			buckets {
				items {
					id
					key
					name
					description
					type
					status
					assetIds
					createdAt
					updatedAt
				}
			}
		}
	`

	var response struct {
		Data struct {
			Buckets struct {
				Items []model.Bucket `json:"items"`
			} `json:"buckets"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	err := s.graphQLClient.ExecuteQueryWithCircuitBreaker(ctx, s.circuitBreaker, s.assetManagerURL+"/graphql", query, &response)

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch buckets from asset-manager", err)
	}

	if err := s.graphQLClient.HandleGraphQLErrors(&response); err != nil {
		return nil, err
	}

	return response.Data.Buckets.Items, nil
}

func (s *Service) fetchAssetFromAssetManager(ctx context.Context, slug string) (*model.Asset, error) {
	query := fmt.Sprintf(`
		query {
			asset(slug: "%s") {
				id
				slug
				title
				description
				type
				genre
				status
				createdAt
				updatedAt
				publishRule {
					isPublic
					publishAt
					unpublishAt
					regions
					ageRating
				}
				videos {
					id
					filename
					status
					format
					url
					thumbnailUrl
					duration
					width
					height
					bitrate
					filesize
				}
			}
		}
	`, slug)

	var response struct {
		Data struct {
			Asset *model.Asset `json:"asset"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	err := s.graphQLClient.ExecuteQueryWithCircuitBreaker(ctx, s.circuitBreaker, s.assetManagerURL+"/graphql", query, &response)

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch asset from asset-manager", err)
	}

	if err := s.graphQLClient.HandleGraphQLErrors(&response); err != nil {
		return nil, err
	}

	return response.Data.Asset, nil
}

func (s *Service) fetchAssetsFromAssetManager(ctx context.Context) ([]model.Asset, error) {
	query := `
		query {
			assets {
				items {
					id
					slug
					title
					description
					type
					genre
					status
					createdAt
					updatedAt
					publishRule {
						isPublic
						publishAt
						unpublishAt
						regions
						ageRating
					}
					videos {
						id
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
						size
						contentType
						status
						thumbnail {
							fileName
							url
							width
							height
							size
							contentType
						}
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	var response struct {
		Data struct {
			Assets struct {
				Items []model.Asset `json:"items"`
			} `json:"assets"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	err := s.graphQLClient.ExecuteQueryWithCircuitBreaker(ctx, s.circuitBreaker, s.assetManagerURL+"/graphql", query, &response)

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch assets from asset-manager", err)
	}

	if err := s.graphQLClient.HandleGraphQLErrors(&response); err != nil {
		return nil, err
	}

	return response.Data.Assets.Items, nil
}

func (s *Service) fetchAssetByIDFromAssetManager(ctx context.Context, id string) (*model.Asset, error) {
	query := fmt.Sprintf(`
		query {
			asset(id: "%s") {
				id
				slug
				title
				description
				type
				genre
				status
				createdAt
				updatedAt
				publishRule {
					isPublic
					publishAt
					unpublishAt
					regions
					ageRating
				}
				videos {
					id
					filename
					status
					format
					url
					thumbnailUrl
					duration
					width
					height
					bitrate
					filesize
				}
			}
		}
	`, id)

	var response struct {
		Data struct {
			Asset *model.Asset `json:"asset"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	err := s.graphQLClient.ExecuteQueryWithCircuitBreaker(ctx, s.circuitBreaker, s.assetManagerURL+"/graphql", query, &response)

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch asset from asset-manager", err)
	}

	if err := s.graphQLClient.HandleGraphQLErrors(&response); err != nil {
		return nil, err
	}

	return response.Data.Asset, nil
}
