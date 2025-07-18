package service

import (
	"context"
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/model"
)

type mockCacheService struct {
	buckets    map[string]*model.Bucket
	assets     map[string]*model.Asset
	allBuckets []model.Bucket
	allAssets  []model.Asset
	shouldErr  bool
}

func newMockCacheService() *mockCacheService {
	return &mockCacheService{
		buckets: make(map[string]*model.Bucket),
		assets:  make(map[string]*model.Asset),
	}
}

func (m *mockCacheService) GetBucket(ctx context.Context, key string) (*model.Bucket, error) {
	if m.shouldErr {
		return nil, errors.NewTransientError("cache error", nil)
	}
	if bucket, exists := m.buckets[key]; exists {
		return bucket, nil
	}
	return nil, nil
}

func (m *mockCacheService) GetBuckets(ctx context.Context) ([]model.Bucket, error) {
	if m.shouldErr {
		return nil, errors.NewTransientError("cache error", nil)
	}
	if m.allBuckets != nil {
		return m.allBuckets, nil
	}
	return nil, nil
}

func (m *mockCacheService) GetAsset(ctx context.Context, slug string) (*model.Asset, error) {
	if m.shouldErr {
		return nil, errors.NewTransientError("cache error", nil)
	}
	if asset, exists := m.assets[slug]; exists {
		return asset, nil
	}
	return nil, nil
}

func (m *mockCacheService) GetAssets(ctx context.Context) ([]model.Asset, error) {
	if m.shouldErr {
		return nil, errors.NewTransientError("cache error", nil)
	}
	if m.allAssets != nil {
		return m.allAssets, nil
	}
	return nil, nil
}

func (m *mockCacheService) SetBucket(ctx context.Context, bucket *model.Bucket) error {
	m.buckets[bucket.Key] = bucket
	return nil
}

func (m *mockCacheService) SetBuckets(ctx context.Context, buckets []model.Bucket) error {
	m.allBuckets = buckets
	return nil
}

func (m *mockCacheService) SetAsset(ctx context.Context, asset *model.Asset) error {
	m.assets[asset.Slug] = asset
	return nil
}

func (m *mockCacheService) SetAssets(ctx context.Context, assets []model.Asset) error {
	m.allAssets = assets
	return nil
}

func (m *mockCacheService) InvalidateBucketCache(ctx context.Context, key string) error {
	return nil
}

func (m *mockCacheService) InvalidateBucketsListCache(ctx context.Context) error {
	return nil
}

func (m *mockCacheService) InvalidateAssetCache(ctx context.Context, slug string) error {
	return nil
}

func (m *mockCacheService) InvalidateAssetsListCache(ctx context.Context) error {
	return nil
}

type mockServiceClient struct{}

func (m *mockServiceClient) GetServiceToken(ctx context.Context) (string, error) {
	return "mock-token", nil
}

func (m *mockServiceClient) GetAuthorizationHeader(ctx context.Context) (string, error) {
	return "Bearer mock-token", nil
}

func stringPtr(s string) *string {
	return &s
}

func createTestService() *Service {
	mockCache := newMockCacheService()
	circuitBreakerConfig := errors.CircuitBreakerConfig{
		Threshold: 5,
		Timeout:   30 * time.Second,
	}

	return NewService(mockCache, "http://localhost:8080", "http://localhost:8081", "test", "test-client", "test-secret", circuitBreakerConfig)
}

func TestService_GetBucket_FromCache(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)

	expectedBucket := &model.Bucket{
		ID:   "test-id",
		Key:  "test-key",
		Name: "Test Bucket",
	}
	mockCache.buckets["test-key"] = expectedBucket

	bucket, err := service.GetBucket(context.Background(), "test-key")

	if err != nil {
		t.Errorf("GetBucket() unexpected error = %v", err)
	}

	if bucket == nil {
		t.Error("GetBucket() expected bucket but got nil")
	}

	if bucket.ID != expectedBucket.ID {
		t.Errorf("GetBucket() bucket ID = %v, want %v", bucket.ID, expectedBucket.ID)
	}
}

func TestService_GetBucket_CacheError(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)
	mockCache.shouldErr = true

	_, err := service.GetBucket(context.Background(), "test-key")

	if err == nil {
		t.Error("GetBucket() expected error but got none")
	}

	if !errors.IsTransient(err) {
		t.Errorf("GetBucket() expected transient error, got %v", err)
	}
}

func TestService_GetBuckets_FromCache(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)

	expectedBuckets := []model.Bucket{
		{ID: "1", Key: "bucket1", Name: "Bucket 1"},
		{ID: "2", Key: "bucket2", Name: "Bucket 2"},
	}
	mockCache.allBuckets = expectedBuckets

	buckets, err := service.GetBuckets(context.Background())

	if err != nil {
		t.Errorf("GetBuckets() unexpected error = %v", err)
	}

	if len(buckets) != len(expectedBuckets) {
		t.Errorf("GetBuckets() returned %d buckets, want %d", len(buckets), len(expectedBuckets))
	}
}

func TestService_GetAsset_FromCache(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)

	expectedAsset := &model.Asset{
		ID:    "test-id",
		Slug:  "test-slug",
		Title: stringPtr("Test Asset"),
	}
	mockCache.assets["test-slug"] = expectedAsset

	asset, err := service.GetAsset(context.Background(), "test-slug")

	if err != nil {
		t.Errorf("GetAsset() unexpected error = %v", err)
	}

	if asset == nil {
		t.Error("GetAsset() expected asset but got nil")
	}

	if asset.ID != expectedAsset.ID {
		t.Errorf("GetAsset() asset ID = %v, want %v", asset.ID, expectedAsset.ID)
	}
}

func TestService_GetAssets_FromCache(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)

	expectedAssets := []model.Asset{
		{ID: "1", Slug: "asset1", Title: stringPtr("Asset 1")},
		{ID: "2", Slug: "asset2", Title: stringPtr("Asset 2")},
	}
	mockCache.allAssets = expectedAssets

	assets, err := service.GetAssets(context.Background())

	if err != nil {
		t.Errorf("GetAssets() unexpected error = %v", err)
	}

	if len(assets) != len(expectedAssets) {
		t.Errorf("GetAssets() returned %d assets, want %d", len(assets), len(expectedAssets))
	}
}

func TestService_GetAssetsInBucket_Success(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)

	bucket := &model.Bucket{
		ID:       "bucket-1",
		Key:      "test-bucket",
		AssetIDs: []string{},
	}
	mockCache.buckets["test-bucket"] = bucket

	assets, err := service.GetAssetsInBucket(context.Background(), "test-bucket")

	if err != nil {
		t.Errorf("GetAssetsInBucket() unexpected error = %v", err)
	}

	if len(assets) != 0 {
		t.Errorf("GetAssetsInBucket() returned %d assets, want 0", len(assets))
	}
}

func TestService_GetAssetsInBucket_BucketNotFound(t *testing.T) {
	service := createTestService()

	_, err := service.GetAssetsInBucket(context.Background(), "non-existent")

	if err == nil {
		t.Error("GetAssetsInBucket() expected error but got none")
	}

	if errors.GetErrorType(err) != errors.ErrorTypeNotFound {
		t.Errorf("GetAssetsInBucket() expected not found error, got %v", err)
	}
}

func TestService_GetAssetsInBucket_EmptyBucket(t *testing.T) {
	service := createTestService()
	mockCache := service.cacheService.(*mockCacheService)

	bucket := &model.Bucket{
		ID:       "bucket-1",
		Key:      "test-bucket",
		AssetIDs: []string{},
	}
	mockCache.buckets["test-bucket"] = bucket

	assets, err := service.GetAssetsInBucket(context.Background(), "test-bucket")

	if err != nil {
		t.Errorf("GetAssetsInBucket() unexpected error = %v", err)
	}

	if len(assets) != 0 {
		t.Errorf("GetAssetsInBucket() returned %d assets, want 0", len(assets))
	}
}
