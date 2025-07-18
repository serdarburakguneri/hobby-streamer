package bucket

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type mockRepository struct {
	buckets map[string]*Bucket
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		buckets: make(map[string]*Bucket),
	}
}

func (m *mockRepository) GetBucketByID(ctx context.Context, id string) (*Bucket, error) {
	if bucket, exists := m.buckets[id]; exists {
		return bucket, nil
	}
	return nil, errors.NewNotFoundError("bucket not found", nil)
}

func (m *mockRepository) GetBucketByKey(ctx context.Context, key string) (*Bucket, error) {
	for _, bucket := range m.buckets {
		if bucket.Key == key {
			return bucket, nil
		}
	}
	return nil, errors.NewNotFoundError("bucket not found", nil)
}

func (m *mockRepository) ListBuckets(ctx context.Context, limit int) (*BucketPage, error) {
	var buckets []Bucket
	count := 0
	for _, bucket := range m.buckets {
		if count >= limit {
			break
		}
		buckets = append(buckets, *bucket)
		count++
	}
	return &BucketPage{
		Items: buckets,
	}, nil
}

func (m *mockRepository) CreateBucket(ctx context.Context, bucket *Bucket) error {
	now := time.Now().UTC()
	if bucket.CreatedAt.IsZero() {
		bucket.CreatedAt = now
	}
	bucket.UpdatedAt = now

	m.buckets[bucket.ID] = bucket
	return nil
}

func (m *mockRepository) UpdateBucket(ctx context.Context, bucket *Bucket) error {
	if _, exists := m.buckets[bucket.ID]; !exists {
		return errors.NewNotFoundError("bucket not found", nil)
	}
	bucket.UpdatedAt = time.Now().UTC()
	m.buckets[bucket.ID] = bucket
	return nil
}

func (m *mockRepository) PatchBucket(ctx context.Context, id string, patch map[string]interface{}) error {
	if _, exists := m.buckets[id]; !exists {
		return errors.NewNotFoundError("bucket not found", nil)
	}
	return nil
}

func (m *mockRepository) DeleteBucket(ctx context.Context, id string) error {
	if _, exists := m.buckets[id]; !exists {
		return errors.NewNotFoundError("bucket not found", nil)
	}
	delete(m.buckets, id)
	return nil
}

func (m *mockRepository) GetBucketsByType(ctx context.Context, bucketType string) ([]Bucket, error) {
	var buckets []Bucket
	for _, bucket := range m.buckets {
		if bucket.Type == bucketType {
			buckets = append(buckets, *bucket)
		}
	}
	return buckets, nil
}

func (m *mockRepository) GetBucketsByAsset(ctx context.Context, assetID string) ([]Bucket, error) {
	var buckets []Bucket
	for _, bucket := range m.buckets {
		for _, bucketAssetID := range bucket.AssetIDs {
			if bucketAssetID == assetID {
				buckets = append(buckets, *bucket)
				break
			}
		}
	}
	return buckets, nil
}

func (m *mockRepository) GetAssetsInBucket(ctx context.Context, bucketID string) ([]string, error) {
	if bucket, exists := m.buckets[bucketID]; exists {
		return bucket.AssetIDs, nil
	}
	return nil, errors.NewNotFoundError("bucket not found", nil)
}

func TestService_CreateBucket(t *testing.T) {
	tests := []struct {
		name    string
		bucket  *Bucket
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid bucket creation",
			bucket: &Bucket{
				Key:         "test-bucket",
				Name:        "Test Bucket",
				Description: "A test bucket",
				Type:        "collection",
			},
			wantErr: false,
		},
		{
			name: "bucket with ID should fail",
			bucket: &Bucket{
				ID:   "123",
				Key:  "test-bucket",
				Name: "Test Bucket",
			},
			wantErr: true,
			errMsg:  "bucket ID should not be set during creation",
		},
		{
			name: "missing key should fail",
			bucket: &Bucket{
				Name: "Test Bucket",
			},
			wantErr: true,
			errMsg:  "key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			service := NewService(repo)
			ctx := context.Background()

			result, err := service.CreateBucket(ctx, tt.bucket)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateBucket() expected error but got none")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("CreateBucket() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateBucket() unexpected error = %v", err)
				return
			}

			if result.ID == "" {
				t.Errorf("CreateBucket() expected ID to be set")
			}

			if result.Status != BucketStatusDraft {
				t.Errorf("CreateBucket() expected status to be %s, got %s", BucketStatusDraft, result.Status)
			}

			if result.CreatedAt.IsZero() {
				t.Errorf("CreateBucket() expected CreatedAt to be set")
			}

			if result.UpdatedAt.IsZero() {
				t.Errorf("CreateBucket() expected UpdatedAt to be set")
			}
		})
	}
}

func TestService_CreateBucket_DuplicateKey(t *testing.T) {
	repo := newMockRepository()
	existingBucket := &Bucket{
		ID:   "existing-id",
		Key:  "test-key",
		Name: "Existing Bucket",
	}
	repo.buckets["existing-id"] = existingBucket

	service := NewService(repo)
	ctx := context.Background()

	newBucket := &Bucket{
		Key:  "test-key",
		Name: "New Bucket",
	}

	_, err := service.CreateBucket(ctx, newBucket)

	if err == nil {
		t.Error("CreateBucket() expected error but got none")
	}

	if !errors.IsConflictError(err) {
		t.Errorf("CreateBucket() expected conflict error, got %v", err)
	}
}

func TestService_UpdateBucket(t *testing.T) {
	repo := newMockRepository()
	existingBucket := &Bucket{
		ID:   "test-id",
		Key:  "test-key",
		Name: "Original Name",
		Type: "collection",
	}
	repo.buckets["test-id"] = existingBucket

	service := NewService(repo)
	ctx := context.Background()

	updatedBucket := &Bucket{
		ID:          "test-id",
		Name:        "Updated Name",
		Description: "Updated description",
		Type:        "playlist",
	}

	err := service.UpdateBucket(ctx, "test-id", updatedBucket)

	if err != nil {
		t.Errorf("UpdateBucket() unexpected error = %v", err)
	}

	updated, err := repo.GetBucketByID(ctx, "test-id")
	if err != nil {
		t.Errorf("Failed to get updated bucket: %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("UpdateBucket() name not updated, got %s", updated.Name)
	}

	if updated.Key != "test-key" {
		t.Errorf("UpdateBucket() key should not change, got %s", updated.Key)
	}
}

func TestService_UpdateBucket_NotFound(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	bucket := &Bucket{
		ID:   "non-existent",
		Name: "Test Bucket",
	}

	err := service.UpdateBucket(ctx, "non-existent", bucket)

	if err == nil {
		t.Error("UpdateBucket() expected error but got none")
	}

	if !errors.IsNotFoundError(err) {
		t.Errorf("UpdateBucket() expected not found error, got %v", err)
	}
}

func TestService_GetBucketByID(t *testing.T) {
	repo := newMockRepository()
	expectedBucket := &Bucket{
		ID:   "test-id",
		Key:  "test-key",
		Name: "Test Bucket",
	}
	repo.buckets["test-id"] = expectedBucket

	service := NewService(repo)
	ctx := context.Background()

	bucket, err := service.GetBucketByID(ctx, "test-id")

	if err != nil {
		t.Errorf("GetBucketByID() unexpected error = %v", err)
	}

	if bucket == nil {
		t.Error("GetBucketByID() expected bucket but got nil")
		return
	}

	if bucket.ID != expectedBucket.ID {
		t.Errorf("GetBucketByID() bucket ID = %v, want %v", bucket.ID, expectedBucket.ID)
	}
}

func TestService_GetBucketByID_NotFound(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	_, err := service.GetBucketByID(ctx, "non-existent")

	if err == nil {
		t.Error("GetBucketByID() expected error but got none")
	}

	if !errors.IsNotFoundError(err) {
		t.Errorf("GetBucketByID() expected not found error, got %v", err)
	}
}

func TestService_GetBucketByKey(t *testing.T) {
	repo := newMockRepository()
	expectedBucket := &Bucket{
		ID:   "test-id",
		Key:  "test-key",
		Name: "Test Bucket",
	}
	repo.buckets["test-id"] = expectedBucket

	service := NewService(repo)
	ctx := context.Background()

	bucket, err := service.GetBucketByKey(ctx, "test-key")

	if err != nil {
		t.Errorf("GetBucketByKey() unexpected error = %v", err)
	}

	if bucket == nil {
		t.Error("GetBucketByKey() expected bucket but got nil")
		return
	}

	if bucket.Key != expectedBucket.Key {
		t.Errorf("GetBucketByKey() bucket key = %v, want %v", bucket.Key, expectedBucket.Key)
	}
}

func TestService_AddAssetToBucket(t *testing.T) {
	repo := newMockRepository()
	bucket := &Bucket{
		ID:       "bucket-1",
		Key:      "test-bucket",
		AssetIDs: []string{"asset-1"},
	}
	repo.buckets["bucket-1"] = bucket

	service := NewService(repo)
	ctx := context.Background()

	err := service.AddAssetToBucket(ctx, "bucket-1", "asset-2")

	if err != nil {
		t.Errorf("AddAssetToBucket() unexpected error = %v", err)
	}

	updatedBucket, err := repo.GetBucketByID(ctx, "bucket-1")
	if err != nil {
		t.Errorf("Failed to get updated bucket: %v", err)
	}

	if len(updatedBucket.AssetIDs) != 2 {
		t.Errorf("AddAssetToBucket() expected 2 assets, got %d", len(updatedBucket.AssetIDs))
	}

	found := false
	for _, assetID := range updatedBucket.AssetIDs {
		if assetID == "asset-2" {
			found = true
			break
		}
	}

	if !found {
		t.Error("AddAssetToBucket() asset-2 not found in bucket")
	}
}

func TestService_AddAssetToBucket_AlreadyExists(t *testing.T) {
	repo := newMockRepository()
	bucket := &Bucket{
		ID:       "bucket-1",
		Key:      "test-bucket",
		AssetIDs: []string{"asset-1"},
	}
	repo.buckets["bucket-1"] = bucket

	service := NewService(repo)
	ctx := context.Background()

	err := service.AddAssetToBucket(ctx, "bucket-1", "asset-1")

	if err == nil {
		t.Error("AddAssetToBucket() expected error but got none")
	}

	if !errors.IsConflictError(err) {
		t.Errorf("AddAssetToBucket() expected conflict error, got %v", err)
	}
}

func TestService_RemoveAssetFromBucket(t *testing.T) {
	repo := newMockRepository()
	bucket := &Bucket{
		ID:       "bucket-1",
		Key:      "test-bucket",
		AssetIDs: []string{"asset-1", "asset-2"},
	}
	repo.buckets["bucket-1"] = bucket

	service := NewService(repo)
	ctx := context.Background()

	err := service.RemoveAssetFromBucket(ctx, "bucket-1", "asset-1")

	if err != nil {
		t.Errorf("RemoveAssetFromBucket() unexpected error = %v", err)
	}

	updatedBucket, err := repo.GetBucketByID(ctx, "bucket-1")
	if err != nil {
		t.Errorf("Failed to get updated bucket: %v", err)
	}

	if len(updatedBucket.AssetIDs) != 1 {
		t.Errorf("RemoveAssetFromBucket() expected 1 asset, got %d", len(updatedBucket.AssetIDs))
	}

	if updatedBucket.AssetIDs[0] != "asset-2" {
		t.Errorf("RemoveAssetFromBucket() expected asset-2, got %s", updatedBucket.AssetIDs[0])
	}
}

func TestService_RemoveAssetFromBucket_NotFound(t *testing.T) {
	repo := newMockRepository()
	bucket := &Bucket{
		ID:       "bucket-1",
		Key:      "test-bucket",
		AssetIDs: []string{"asset-1"},
	}
	repo.buckets["bucket-1"] = bucket

	service := NewService(repo)
	ctx := context.Background()

	err := service.RemoveAssetFromBucket(ctx, "bucket-1", "asset-2")

	if err == nil {
		t.Error("RemoveAssetFromBucket() expected error but got none")
	}

	if !errors.IsNotFoundError(err) {
		t.Errorf("RemoveAssetFromBucket() expected not found error, got %v", err)
	}
}

func TestService_GetBucketsByType(t *testing.T) {
	repo := newMockRepository()
	bucket1 := &Bucket{ID: "1", Key: "bucket1", Type: "collection"}
	bucket2 := &Bucket{ID: "2", Key: "bucket2", Type: "playlist"}
	bucket3 := &Bucket{ID: "3", Key: "bucket3", Type: "collection"}
	repo.buckets["1"] = bucket1
	repo.buckets["2"] = bucket2
	repo.buckets["3"] = bucket3

	service := NewService(repo)
	ctx := context.Background()

	buckets, err := service.GetBucketsByType(ctx, "collection")

	if err != nil {
		t.Errorf("GetBucketsByType() unexpected error = %v", err)
	}

	if len(buckets) != 2 {
		t.Errorf("GetBucketsByType() expected 2 buckets, got %d", len(buckets))
	}

	for _, bucket := range buckets {
		if bucket.Type != "collection" {
			t.Errorf("GetBucketsByType() expected type 'collection', got %s", bucket.Type)
		}
	}
}

func TestService_GetBucketsByType_EmptyType(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	_, err := service.GetBucketsByType(ctx, "")

	if err == nil {
		t.Error("GetBucketsByType() expected error but got none")
	}

	if !errors.IsValidationError(err) {
		t.Errorf("GetBucketsByType() expected validation error, got %v", err)
	}
}
