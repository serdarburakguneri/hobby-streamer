package bucket

import (
	"context"
	"testing"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, bucket *Bucket) error {
	args := m.Called(ctx, bucket)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*Bucket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Bucket), args.Error(1)
}

func (m *MockRepository) GetByKey(ctx context.Context, key string) (*Bucket, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Bucket), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, bucket *Bucket) error {
	args := m.Called(ctx, bucket)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, limit *int, lastKey map[string]interface{}) (*BucketPage, error) {
	args := m.Called(ctx, limit, lastKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BucketPage), args.Error(1)
}

func (m *MockRepository) Search(ctx context.Context, query string, limit *int, lastKey map[string]interface{}) (*BucketPage, error) {
	args := m.Called(ctx, query, limit, lastKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BucketPage), args.Error(1)
}

func (m *MockRepository) GetByOwnerID(ctx context.Context, ownerID string, limit *int, lastKey map[string]interface{}) (*BucketPage, error) {
	args := m.Called(ctx, ownerID, limit, lastKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BucketPage), args.Error(1)
}

func (m *MockRepository) AddAsset(ctx context.Context, bucketID string, assetID string) error {
	args := m.Called(ctx, bucketID, assetID)
	return args.Error(0)
}

func (m *MockRepository) RemoveAsset(ctx context.Context, bucketID string, assetID string) error {
	args := m.Called(ctx, bucketID, assetID)
	return args.Error(0)
}

func (m *MockRepository) GetAssetIDs(ctx context.Context, bucketID string, limit *int, lastKey map[string]interface{}) ([]string, error) {
	args := m.Called(ctx, bucketID, limit, lastKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) HasAsset(ctx context.Context, bucketID string, assetID string) (bool, error) {
	args := m.Called(ctx, bucketID, assetID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) AssetCount(ctx context.Context, bucketID string) (int, error) {
	args := m.Called(ctx, bucketID)
	return args.Int(0), args.Error(1)
}

func TestDomainService_IsKeyAvailable(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewDomainService(mockRepo)

	tests := []struct {
		name        string
		key         string
		setupMock   func()
		expected    bool
		expectError bool
	}{
		{
			name: "key available",
			key:  "available-key",
			setupMock: func() {
				mockRepo.On("GetByKey", ctx, "available-key").Return(nil, pkgerrors.NewNotFoundError("bucket not found", nil))
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "key not available",
			key:  "existing-key",
			setupMock: func() {
				bucket, _ := NewBucket("Existing Bucket", "existing-key")
				mockRepo.On("GetByKey", ctx, "existing-key").Return(bucket, nil)
			},
			expected:    false,
			expectError: false,
		},
		{
			name: "repository error",
			key:  "error-key",
			setupMock: func() {
				mockRepo.On("GetByKey", ctx, "error-key").Return(nil, pkgerrors.NewInternalError("database error", nil))
			},
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			result, err := service.IsKeyAvailable(ctx, tt.key)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDomainService_ValidateBucketOwnership(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewDomainService(mockRepo)

	ownerID := "user123"
	bucket := ReconstructBucket(
		"bucket123",
		"Test Bucket",
		nil,
		"test-bucket",
		&ownerID,
		nil,
		nil,
		time.Now(),
		time.Now(),
	)

	tests := []struct {
		name        string
		bucketID    string
		userID      string
		setupMock   func()
		expectError bool
	}{
		{
			name:     "valid ownership",
			bucketID: "bucket123",
			userID:   "user123",
			setupMock: func() {
				mockRepo.On("GetByID", ctx, "bucket123").Return(bucket, nil)
			},
			expectError: false,
		},
		{
			name:     "invalid ownership",
			bucketID: "bucket123",
			userID:   "user456",
			setupMock: func() {
				mockRepo.On("GetByID", ctx, "bucket123").Return(bucket, nil)
			},
			expectError: true,
		},
		{
			name:     "bucket not found",
			bucketID: "nonexistent",
			userID:   "user123",
			setupMock: func() {
				mockRepo.On("GetByID", ctx, "nonexistent").Return(nil, pkgerrors.NewNotFoundError("bucket not found", nil))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			err := service.ValidateBucketOwnership(ctx, tt.bucketID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDomainService_ValidateBucketNotEmpty(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewDomainService(mockRepo)

	tests := []struct {
		name        string
		bucketID    string
		setupMock   func()
		expectError bool
	}{
		{
			name:     "bucket has assets",
			bucketID: "bucket123",
			setupMock: func() {
				mockRepo.On("AssetCount", ctx, "bucket123").Return(5, nil)
			},
			expectError: false,
		},
		{
			name:     "bucket is empty",
			bucketID: "bucket123",
			setupMock: func() {
				mockRepo.On("AssetCount", ctx, "bucket123").Return(0, nil)
			},
			expectError: true,
		},
		{
			name:     "repository error",
			bucketID: "bucket123",
			setupMock: func() {
				mockRepo.On("AssetCount", ctx, "bucket123").Return(0, pkgerrors.NewInternalError("database error", nil))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			err := service.ValidateBucketNotEmpty(ctx, tt.bucketID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDomainService_ValidateBucketOwnership_NoOwner(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewDomainService(mockRepo)

	bucket := ReconstructBucket(
		"bucket123",
		"Test Bucket",
		nil,
		"test-bucket",
		nil,
		nil,
		nil,
		time.Now(),
		time.Now(),
	)

	mockRepo.On("GetByID", ctx, "bucket123").Return(bucket, nil)

	err := service.ValidateBucketOwnership(ctx, "bucket123", "user123")
	assert.Error(t, err)
	assert.Equal(t, pkgerrors.ErrorTypeForbidden, pkgerrors.GetErrorType(err))

	mockRepo.AssertExpectations(t)
}
