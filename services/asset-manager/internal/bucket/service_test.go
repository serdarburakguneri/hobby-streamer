package bucket_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/bucket"
)

type MockRepository struct {
	SaveBucketFunc    func(ctx context.Context, b *bucket.Bucket) error
	GetBucketByIDFunc func(ctx context.Context, id int) (*bucket.Bucket, error)
}

func (m *MockRepository) SaveBucket(ctx context.Context, b *bucket.Bucket) error {
	if m.SaveBucketFunc != nil {
		return m.SaveBucketFunc(ctx, b)
	}
	return nil
}
func (m *MockRepository) GetBucketByID(ctx context.Context, id int) (*bucket.Bucket, error) {
	if m.GetBucketByIDFunc != nil {
		return m.GetBucketByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) AddAssetToBucket(ctx context.Context, id int, assetID int) error {
	return nil
}
func (m *MockRepository) RemoveAssetFromBucket(ctx context.Context, id int, assetID int) error {
	return nil
}
func (m *MockRepository) ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*bucket.BucketPage, error) {
	return nil, nil
}
func (m *MockRepository) UpdateBucket(ctx context.Context, id int, b *bucket.Bucket) error {
	return nil
}
func (m *MockRepository) PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error {
	return nil
}

func (m *MockRepository) CreateBucket(ctx context.Context, b *bucket.Bucket) (*bucket.Bucket, error) {
	return nil, nil
}

var _ bucket.BucketService = (*MockRepository)(nil)

func TestCreateBucket_Success(t *testing.T) {
	mockRepo := &MockRepository{
		SaveBucketFunc: func(ctx context.Context, b *bucket.Bucket) error {
			return nil
		},
	}
	svc := bucket.NewService(mockRepo)

	bucketObj := &bucket.Bucket{}
	got, err := svc.CreateBucket(context.Background(), bucketObj)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID == "" {
		t.Errorf("expected bucket ID to be set, got empty string")
	}
}

func TestCreateBucket_WithIDSet(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := bucket.NewService(mockRepo)

	bucketObj := &bucket.Bucket{ID: "123"}
	_, err := svc.CreateBucket(context.Background(), bucketObj)
	if err == nil {
		t.Fatalf("expected error when ID is set, got nil")
	}
}

func TestAddAssetToBucket_Duplicate(t *testing.T) {
	mockRepo := &MockRepository{
		GetBucketByIDFunc: func(ctx context.Context, id int) (*bucket.Bucket, error) {
			return &bucket.Bucket{
				AssetIDs: []string{"42"},
			}, nil
		},
		SaveBucketFunc: func(ctx context.Context, b *bucket.Bucket) error {
			return nil
		},
	}
	svc := bucket.NewService(mockRepo)
	assetID := 42
	err := svc.AddAssetToBucket(context.Background(), 1, assetID)
	if err == nil || err.Error() != "asset already in bucket" {
		t.Errorf("expected duplicate error, got %v", err)
	}
}
