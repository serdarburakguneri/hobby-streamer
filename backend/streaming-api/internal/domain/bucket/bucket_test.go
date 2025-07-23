package bucket

import (
	"testing"
	"time"

	assetdomain "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	"github.com/stretchr/testify/assert"
)

func TestBucket_IsActive(t *testing.T) {
	b := &Bucket{status: nil}
	assert.True(t, b.IsActive())
	active := BucketStatus{"active"}
	b2 := &Bucket{status: &active}
	assert.True(t, b2.IsActive())
	inactive := BucketStatus{"inactive"}
	b3 := &Bucket{status: &inactive}
	assert.False(t, b3.IsActive())
}

func TestBucket_IsCollection(t *testing.T) {
	col := BucketType{"collection"}
	b := &Bucket{bucketType: col}
	assert.True(t, b.IsCollection())
}

func TestBucket_IsPlaylist(t *testing.T) {
	pl := BucketType{"playlist"}
	b := &Bucket{bucketType: pl}
	assert.True(t, b.IsPlaylist())
}

func TestBucket_IsCategory(t *testing.T) {
	cat := BucketType{"category"}
	b := &Bucket{bucketType: cat}
	assert.True(t, b.IsCategory())
}

func TestBucket_IsFeatured(t *testing.T) {
	ft := BucketType{"featured"}
	b := &Bucket{bucketType: ft}
	assert.True(t, b.IsFeatured())
}

func TestBucket_IsTrending(t *testing.T) {
	tr := BucketType{"trending"}
	b := &Bucket{bucketType: tr}
	assert.True(t, b.IsTrending())
}

func TestBucket_ContainsAsset(t *testing.T) {
	ids, _ := NewAssetIDs([]string{"a1", "a2"})
	b := &Bucket{assetIDs: ids}
	assert.True(t, b.ContainsAsset("a1"))
	assert.False(t, b.ContainsAsset("a3"))
}

func TestBucket_GetPublicAssets(t *testing.T) {
	now := time.Now().UTC()
	publishAt := now.Add(-time.Hour)
	unpublishAt := now.Add(time.Hour)
	pr, _ := assetdomain.NewPublishRuleValue(&publishAt, &unpublishAt, nil, nil)
	id, _ := assetdomain.NewAssetID("id")
	slug, _ := assetdomain.NewSlug("slug")
	typeVO, _ := assetdomain.NewAssetType("movie")
	createdAt := assetdomain.NewCreatedAt(now)
	updatedAt := assetdomain.NewUpdatedAt(now)
	pub := assetdomain.NewAsset(
		*id,
		*slug,
		nil, nil,
		*typeVO,
		nil, nil, nil, nil,
		*createdAt,
		*updatedAt,
		nil, nil, nil, nil, pr,
	)
	b := &Bucket{assets: []*assetdomain.Asset{pub}}
	assert.Len(t, b.GetPublicAssets(), 1)
}
