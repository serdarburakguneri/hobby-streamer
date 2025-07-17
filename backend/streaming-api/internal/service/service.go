package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"

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
	cacheService    *cache.Service
	logger          *logger.Logger
	assetManagerURL string
	serviceClient   *auth.ServiceClient
	circuitBreaker  *errors.CircuitBreaker
}

func NewService(cacheService *cache.Service, assetManagerURL, keycloakURL, realm, clientID, clientSecret string) *Service {
	log := logger.Get().WithService("streaming-service")
	log.Info("Initializing streaming service", "asset_manager_url", assetManagerURL, "keycloak_url", keycloakURL, "realm", realm, "client_id", clientID)

	serviceClient := auth.NewServiceClient(keycloakURL, realm, clientID, clientSecret)
	log.Info("Service client created successfully")

	circuitBreaker := errors.NewCircuitBreaker(errors.CircuitBreakerConfig{
		Name:      "asset-manager",
		Threshold: 5,
		Timeout:   30 * time.Second,
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

	var assets []model.Asset
	for _, assetID := range bucket.AssetIDs {
		asset, err := s.fetchAssetByIDFromAssetManager(ctx, assetID)
		if err != nil {
			s.logger.WithError(err).Warn("Failed to fetch asset", "asset_id", assetID)
			continue
		}
		if asset != nil {
			assets = append(assets, *asset)
		}
	}

	return assets, nil
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
			Bucket *model.Bucket `json:"bucket"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	err := s.circuitBreaker.Execute(ctx, func() error {
		return s.makeGraphQLRequest(ctx, s.assetManagerURL+"/graphql", query, &response)
	})

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch bucket from asset-manager", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.NewExternalError(fmt.Sprintf("GraphQL errors: %v", response.Errors), nil)
	}

	return response.Data.Bucket, nil
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

	err := s.circuitBreaker.Execute(ctx, func() error {
		return s.makeGraphQLRequest(ctx, s.assetManagerURL+"/graphql", query, &response)
	})

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch buckets from asset-manager", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.NewExternalError(fmt.Sprintf("GraphQL errors: %v", response.Errors), nil)
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
				bucketKey
				createdAt
				updatedAt
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

	err := s.circuitBreaker.Execute(ctx, func() error {
		return s.makeGraphQLRequest(ctx, s.assetManagerURL+"/graphql", query, &response)
	})

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch asset from asset-manager", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.NewExternalError(fmt.Sprintf("GraphQL errors: %v", response.Errors), nil)
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

	err := s.circuitBreaker.Execute(ctx, func() error {
		return s.makeGraphQLRequest(ctx, s.assetManagerURL+"/graphql", query, &response)
	})

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch assets from asset-manager", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.NewExternalError(fmt.Sprintf("GraphQL errors: %v", response.Errors), nil)
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

	err := s.circuitBreaker.Execute(ctx, func() error {
		return s.makeGraphQLRequest(ctx, s.assetManagerURL+"/graphql", query, &response)
	})

	if err != nil {
		if errors.IsAppError(err) && errors.GetErrorType(err) == errors.ErrorTypeCircuitBreaker {
			return nil, errors.NewExternalError("asset-manager service unavailable", err)
		}
		return nil, errors.NewExternalError("failed to fetch asset from asset-manager", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.NewExternalError(fmt.Sprintf("GraphQL errors: %v", response.Errors), nil)
	}

	return response.Data.Asset, nil
}

func (s *Service) makeGraphQLRequest(ctx context.Context, url, query string, response interface{}) error {
	requestBody := map[string]string{
		"query": query,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return errors.NewInternalError("failed to marshal request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return errors.NewInternalError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	s.logger.Info("Getting service token for request", "url", url)
	authHeader, err := s.serviceClient.GetAuthorizationHeader(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get service token")
		return errors.NewExternalError("failed to get service token", err)
	}

	s.logger.Info("Service token obtained successfully", "auth_header_length", len(authHeader))
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.NewTransientError("failed to make request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("GraphQL request failed", "status_code", resp.StatusCode, "url", url)
		return errors.NewExternalError(fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return errors.NewInternalError("failed to decode response", err)
	}

	return nil
}
