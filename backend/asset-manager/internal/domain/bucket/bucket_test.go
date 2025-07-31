package bucket

import (
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestNewBucket(t *testing.T) {
	tests := []struct {
		name        string
		bucketName  string
		bucketKey   string
		expectError bool
	}{
		{
			name:        "valid bucket",
			bucketName:  "My Bucket",
			bucketKey:   "my-bucket",
			expectError: false,
		},
		{
			name:        "empty name",
			bucketName:  "",
			bucketKey:   "my-bucket",
			expectError: true,
		},
		{
			name:        "invalid key format",
			bucketName:  "My Bucket",
			bucketKey:   "MY_BUCKET",
			expectError: true,
		},
		{
			name:        "key too short",
			bucketName:  "My Bucket",
			bucketKey:   "ab",
			expectError: true,
		},
		{
			name:        "key too long",
			bucketName:  "My Bucket",
			bucketKey:   "this-is-a-very-long-bucket-key-that-exceeds-the-maximum-length",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket, err := entity.NewBucket(tt.bucketName, tt.bucketKey)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, bucket)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bucket)
				assert.Equal(t, tt.bucketName, bucket.Name().Value())
				assert.Equal(t, tt.bucketKey, bucket.Key().Value())
				assert.NotEmpty(t, bucket.ID().Value())
				assert.NotZero(t, bucket.CreatedAt().Value())
				assert.NotZero(t, bucket.UpdatedAt().Value())
			}
		})
	}
}

func TestReconstructBucket(t *testing.T) {
	now := time.Now().UTC()
	metadata := map[string]interface{}{"key": "value"}
	bucketID, _ := valueobjects.NewBucketID("bucket123")
	bucketName, _ := valueobjects.NewBucketName("Test Bucket")
	bucketKey, _ := valueobjects.NewBucketKey("test-bucket")
	ownerID, _ := valueobjects.NewOwnerID("user123")
	createdAt := valueobjects.NewCreatedAt(now)
	updatedAt := valueobjects.NewUpdatedAt(now)

	bucket := entity.ReconstructBucket(
		*bucketID,
		*bucketName,
		nil,
		*bucketKey,
		nil,
		ownerID,
		nil,
		metadata,
		createdAt,
		updatedAt,
		[]string{},
	)

	assert.Equal(t, "bucket123", bucket.ID().Value())
	assert.Equal(t, "Test Bucket", bucket.Name().Value())
	assert.Equal(t, "test-bucket", bucket.Key().Value())
	assert.Equal(t, "user123", bucket.OwnerID().Value())
	assert.Equal(t, metadata, bucket.Metadata())
	assert.Equal(t, now, bucket.CreatedAt().Value())
	assert.Equal(t, now, bucket.UpdatedAt().Value())
}
