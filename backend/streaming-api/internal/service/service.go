package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

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
}

func NewService(cacheService *cache.Service) *Service {
	assetManagerURL := getEnv("ASSET_MANAGER_URL", "http://localhost:8081")
	return &Service{
		cacheService:    cacheService,
		logger:          logger.Get().WithService("streaming-service"),
		assetManagerURL: assetManagerURL,
	}
}

func (s *Service) GetBucket(ctx context.Context, key string) (*model.Bucket, error) {
	bucket, err := s.cacheService.GetBucket(ctx, key)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get bucket from cache", "key", key)
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if bucket != nil {
		s.logger.Debug("Bucket found in cache", "key", key)
		return bucket, nil
	}

	s.logger.Debug("Bucket not found in cache, fetching from asset-manager", "key", key)

	bucket, err = s.fetchBucketFromAssetManager(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bucket from asset-manager: %w", err)
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
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if buckets != nil {
		s.logger.Debug("Buckets found in cache", "count", len(buckets))
		return buckets, nil
	}

	s.logger.Debug("Buckets not found in cache, fetching from asset-manager")

	buckets, err = s.fetchBucketsFromAssetManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch buckets from asset-manager: %w", err)
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
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if asset != nil {
		s.logger.Debug("Asset found in cache", "slug", slug)
		return asset, nil
	}

	s.logger.Debug("Asset not found in cache, fetching from asset-manager", "slug", slug)

	asset, err = s.fetchAssetFromAssetManager(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset from asset-manager: %w", err)
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
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if assets != nil {
		s.logger.Debug("Assets found in cache", "count", len(assets))
		return assets, nil
	}

	s.logger.Debug("Assets not found in cache, fetching from asset-manager")

	assets, err = s.fetchAssetsFromAssetManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assets from asset-manager: %w", err)
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
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	if bucket == nil {
		return nil, fmt.Errorf("bucket not found: %s", bucketKey)
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
	url := fmt.Sprintf("%s/graphql", s.assetManagerURL)

	query := `{
		"query": "query GetBucket($key: String!) { bucketByKey(key: $key) { id key name description type status assetIds createdAt updatedAt } }",
		"variables": {"key": "` + key + `"}
	}`

	var response struct {
		Data struct {
			BucketByKey *model.Bucket `json:"bucketByKey"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	if err := s.makeGraphQLRequest(ctx, url, query, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", response.Errors)
	}

	return response.Data.BucketByKey, nil
}

func (s *Service) fetchBucketsFromAssetManager(ctx context.Context) ([]model.Bucket, error) {
	url := fmt.Sprintf("%s/graphql", s.assetManagerURL)

	query := `{
		"query": "query GetBuckets { buckets { items { id key name description type status assetIds createdAt updatedAt } } }"
	}`

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

	if err := s.makeGraphQLRequest(ctx, url, query, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", response.Errors)
	}

	return response.Data.Buckets.Items, nil
}

func (s *Service) fetchAssetFromAssetManager(ctx context.Context, slug string) (*model.Asset, error) {
	url := fmt.Sprintf("%s/graphql", s.assetManagerURL)

	query := `{
		"query": "query GetAsset($slug: String!) { assetBySlug(slug: $slug) { id slug title description type genre genres tags status createdAt updatedAt metadata ownerId videos { id type format storageLocation { bucket key url } width height duration bitrate codec size contentType streamInfo { downloadUrl cdnPrefix playUrl } metadata status thumbnail { fileName url storageLocation { bucket key url } width height size contentType metadata } createdAt updatedAt } publishRule { isPublic publishAt unpublishAt regions ageRating } } }",
		"variables": {"slug": "` + slug + `"}
	}`

	var response struct {
		Data struct {
			AssetBySlug *model.Asset `json:"assetBySlug"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	if err := s.makeGraphQLRequest(ctx, url, query, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", response.Errors)
	}

	return response.Data.AssetBySlug, nil
}

func (s *Service) fetchAssetByIDFromAssetManager(ctx context.Context, id string) (*model.Asset, error) {
	url := fmt.Sprintf("%s/graphql", s.assetManagerURL)

	query := `{
		"query": "query GetAsset($id: ID!) { asset(id: $id) { id slug title description type genre genres tags status createdAt updatedAt metadata ownerId videos { id type format storageLocation { bucket key url } width height duration bitrate codec size contentType streamInfo { downloadUrl cdnPrefix playUrl } metadata status thumbnail { fileName url storageLocation { bucket key url } width height size contentType metadata } createdAt updatedAt } publishRule { isPublic publishAt unpublishAt regions ageRating } } }",
		"variables": {"id": "` + id + `"}
	}`

	var response struct {
		Data struct {
			Asset *model.Asset `json:"asset"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	if err := s.makeGraphQLRequest(ctx, url, query, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", response.Errors)
	}

	return response.Data.Asset, nil
}

func (s *Service) fetchAssetsFromAssetManager(ctx context.Context) ([]model.Asset, error) {
	url := fmt.Sprintf("%s/graphql", s.assetManagerURL)

	query := `{
		"query": "query GetAssets { assets { items { id slug title description type genre genres tags status createdAt updatedAt metadata ownerId videos { id type format storageLocation { bucket key url } width height duration bitrate codec size contentType streamInfo { downloadUrl cdnPrefix playUrl } metadata status thumbnail { fileName url storageLocation { bucket key url } width height size contentType metadata } createdAt updatedAt } publishRule { isPublic publishAt unpublishAt regions ageRating } } } }"
	}`

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

	if err := s.makeGraphQLRequest(ctx, url, query, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", response.Errors)
	}

	return response.Data.Assets.Items, nil
}

func (s *Service) makeGraphQLRequest(ctx context.Context, url, query string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(query))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
