package bucket

import (
	"testing"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
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
			bucket, err := NewBucket(tt.bucketName, tt.bucketKey)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, bucket)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bucket)
				assert.Equal(t, tt.bucketName, bucket.Name())
				assert.Equal(t, tt.bucketKey, bucket.Key())
				assert.NotEmpty(t, bucket.ID())
				assert.NotZero(t, bucket.CreatedAt())
				assert.NotZero(t, bucket.UpdatedAt())
			}
		})
	}
}

func TestReconstructBucket(t *testing.T) {
	now := time.Now().UTC()
	description := "Test description"
	ownerID := "user123"
	status := "active"
	metadata := map[string]interface{}{"key": "value"}

	bucket := ReconstructBucket(
		"bucket123",
		"Test Bucket",
		&description,
		"test-bucket",
		&ownerID,
		&status,
		metadata,
		now,
		now,
	)

	assert.Equal(t, "bucket123", bucket.ID())
	assert.Equal(t, "Test Bucket", bucket.Name())
	assert.Equal(t, &description, bucket.Description())
	assert.Equal(t, "test-bucket", bucket.Key())
	assert.Equal(t, &ownerID, bucket.OwnerID())
	assert.Equal(t, &status, bucket.Status())
	assert.Equal(t, metadata, bucket.Metadata())
	assert.Equal(t, now, bucket.CreatedAt())
	assert.Equal(t, now, bucket.UpdatedAt())
}

func TestBucket_UpdateName(t *testing.T) {
	tests := []struct {
		name         string
		originalName string
		newName      string
		expectError  bool
	}{
		{
			name:         "valid name",
			originalName: "Original Name",
			newName:      "New Name",
			expectError:  false,
		},
		{
			name:         "empty name",
			originalName: "Original Name",
			newName:      "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket, _ := NewBucket(tt.originalName, "test-bucket")
			err := bucket.UpdateName(tt.newName)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.originalName, bucket.Name())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, bucket.Name())
			}
		})
	}
}

func TestBucket_ValidateName(t *testing.T) {
	bucket, _ := NewBucket("Test Bucket", "test-bucket")

	tests := []struct {
		name        string
		bucketName  string
		expectError bool
	}{
		{
			name:        "valid name",
			bucketName:  "Valid Name",
			expectError: false,
		},
		{
			name:        "empty name",
			bucketName:  "",
			expectError: true,
		},
		{
			name:        "name too long",
			bucketName:  "This is a very long bucket name that exceeds the maximum allowed length of one hundred characters and should cause a validation error",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket.name = tt.bucketName
			err := bucket.ValidateName()
			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, pkgerrors.IsValidationError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucket_ValidateKey(t *testing.T) {
	bucket, _ := NewBucket("Test Bucket", "test-bucket")

	tests := []struct {
		name        string
		bucketKey   string
		expectError bool
	}{
		{
			name:        "valid key",
			bucketKey:   "valid-key",
			expectError: false,
		},
		{
			name:        "key with uppercase",
			bucketKey:   "Invalid-Key",
			expectError: true,
		},
		{
			name:        "key with underscore",
			bucketKey:   "invalid_key",
			expectError: true,
		},
		{
			name:        "key too short",
			bucketKey:   "ab",
			expectError: true,
		},
		{
			name:        "key too long",
			bucketKey:   "this-is-a-very-long-bucket-key-that-exceeds-the-maximum-length",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket.key = tt.bucketKey
			err := bucket.ValidateKey()
			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, pkgerrors.IsValidationError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucket_ValidateOwnership(t *testing.T) {
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
		userID      string
		expectError bool
	}{
		{
			name:        "valid owner",
			userID:      "user123",
			expectError: false,
		},
		{
			name:        "wrong owner",
			userID:      "user456",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bucket.ValidateOwnership(tt.userID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, pkgerrors.ErrorTypeForbidden, pkgerrors.GetErrorType(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucket_ValidateOwnership_NoOwner(t *testing.T) {
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

	err := bucket.ValidateOwnership("user123")
	assert.Error(t, err)
	assert.Equal(t, pkgerrors.ErrorTypeForbidden, pkgerrors.GetErrorType(err))
}

func TestBucket_CanAddAsset(t *testing.T) {
	bucket, _ := NewBucket("Test Bucket", "test-bucket")

	tests := []struct {
		name        string
		assetID     string
		hasAsset    bool
		expectError bool
	}{
		{
			name:        "asset not exists",
			assetID:     "asset123",
			hasAsset:    false,
			expectError: false,
		},
		{
			name:        "asset already exists",
			assetID:     "asset123",
			hasAsset:    true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAssetFunc := func(bucketID, assetID string) (bool, error) {
				return tt.hasAsset, nil
			}

			err := bucket.CanAddAsset(tt.assetID, hasAssetFunc)
			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, pkgerrors.IsValidationError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucket_ValidateNotEmpty(t *testing.T) {
	bucket, _ := NewBucket("Test Bucket", "test-bucket")

	tests := []struct {
		name        string
		assetCount  int
		expectError bool
	}{
		{
			name:        "bucket has assets",
			assetCount:  5,
			expectError: false,
		},
		{
			name:        "bucket is empty",
			assetCount:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assetCountFunc := func(bucketID string) (int, error) {
				return tt.assetCount, nil
			}

			err := bucket.ValidateNotEmpty(assetCountFunc)
			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, pkgerrors.IsValidationError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucket_Setters(t *testing.T) {
	bucket, _ := NewBucket("Test Bucket", "test-bucket")
	originalUpdatedAt := bucket.UpdatedAt()

	time.Sleep(1 * time.Millisecond)

	description := "New description"
	bucket.UpdateDescription(&description)
	assert.Equal(t, &description, bucket.Description())
	assert.True(t, bucket.UpdatedAt().After(originalUpdatedAt))

	time.Sleep(1 * time.Millisecond)
	originalUpdatedAt = bucket.UpdatedAt()

	ownerID := "user123"
	bucket.SetOwnerID(&ownerID)
	assert.Equal(t, &ownerID, bucket.OwnerID())
	assert.True(t, bucket.UpdatedAt().After(originalUpdatedAt))

	time.Sleep(1 * time.Millisecond)
	originalUpdatedAt = bucket.UpdatedAt()

	metadata := map[string]interface{}{"key": "value"}
	bucket.SetMetadata(metadata)
	assert.Equal(t, metadata, bucket.Metadata())
	assert.True(t, bucket.UpdatedAt().After(originalUpdatedAt))

	time.Sleep(1 * time.Millisecond)
	originalUpdatedAt = bucket.UpdatedAt()

	status := "active"
	bucket.SetStatus(&status)
	assert.Equal(t, &status, bucket.Status())
	assert.True(t, bucket.UpdatedAt().After(originalUpdatedAt))
}

func TestIsValidKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"valid key", "valid-key", true},
		{"valid key with numbers", "key123", true},
		{"valid key with hyphens", "my-bucket-key", true},
		{"uppercase letters", "Invalid-Key", false},
		{"underscore", "invalid_key", false},
		{"spaces", "invalid key", false},
		{"special characters", "key@bucket", false},
		{"too short", "ab", false},
		{"too long", "this-is-a-very-long-bucket-key-that-exceeds-the-maximum-length", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}
