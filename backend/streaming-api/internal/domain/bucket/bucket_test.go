package bucket

import (
	"testing"
	"time"

	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
	bucketvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestBucket_IsActive(t *testing.T) {
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("collection")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b.IsActive())

	active, _ := bucketvalueobjects.NewBucketStatus("active")
	b2 := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, active, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b2.IsActive())

	inactive, _ := bucketvalueobjects.NewBucketStatus("inactive")
	b3 := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, inactive, assetIDs, *createdAt, *updatedAt, nil)
	assert.False(t, b3.IsActive())
}

func TestBucket_IsCollection(t *testing.T) {
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("collection")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b.IsCollection())
}

func TestBucket_IsPlaylist(t *testing.T) {
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("playlist")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b.IsPlaylist())
}

func TestBucket_IsCategory(t *testing.T) {
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("category")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b.IsCategory())
}

func TestBucket_IsFeatured(t *testing.T) {
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("featured")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b.IsFeatured())
}

func TestBucket_IsTrending(t *testing.T) {
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("trending")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *createdAt, *updatedAt, nil)
	assert.True(t, b.IsTrending())
}

func TestBucket_ContainsAsset(t *testing.T) {
	ids, _ := bucketvalueobjects.NewAssetIDs([]string{"a1", "a2"})
	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("collection")
	createdAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, ids, *createdAt, *updatedAt, nil)
	assert.True(t, b.ContainsAsset("a1"))
	assert.False(t, b.ContainsAsset("a3"))
}

func TestBucket_GetPublicAssets(t *testing.T) {
	now := time.Now().UTC()
	publishAt := now.Add(-time.Hour)
	unpublishAt := now.Add(time.Hour)
	pr, _ := assetvalueobjects.NewPublishRuleValue(&publishAt, &unpublishAt, nil, nil)
	id, _ := assetvalueobjects.NewAssetID("id")
	slug, _ := assetvalueobjects.NewSlug("slug")
	typeVO, _ := assetvalueobjects.NewAssetType("movie")
	createdAt := assetvalueobjects.NewCreatedAt(now)
	updatedAt := assetvalueobjects.NewUpdatedAt(now)
	pub := assetentity.NewAsset(
		*id,
		*slug,
		nil, nil,
		*typeVO,
		nil, nil, nil, nil,
		createdAt.Value(),
		updatedAt.Value(),
		nil, nil, nil, nil, pr,
	)

	bucketID, _ := bucketvalueobjects.NewBucketID("test-id")
	bucketKey, _ := bucketvalueobjects.NewBucketKey("test-key")
	bucketName, _ := bucketvalueobjects.NewBucketName("test-name")
	bucketType, _ := bucketvalueobjects.NewBucketType("collection")
	bucketCreatedAt := bucketvalueobjects.NewCreatedAt(time.Now().UTC())
	bucketUpdatedAt := bucketvalueobjects.NewUpdatedAt(time.Now().UTC())
	assetIDs, _ := bucketvalueobjects.NewAssetIDs([]string{})

	b := bucketentity.NewBucket(*bucketID, *bucketKey, *bucketName, nil, *bucketType, nil, assetIDs, *bucketCreatedAt, *bucketUpdatedAt, []*assetentity.Asset{pub})
	assert.Len(t, b.GetPublicAssets(), 1)
}
