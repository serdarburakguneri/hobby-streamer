package asset

import (
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestValueObjects(t *testing.T) {
	t.Run("AssetID", func(t *testing.T) {
		// Valid AssetID
		assetID, err := NewAssetID("asset-123")
		assert.NoError(t, err)
		assert.Equal(t, "asset-123", assetID.Value())

		// Invalid AssetID
		_, err = NewAssetID("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAssetID, err)
	})

	t.Run("Slug", func(t *testing.T) {
		// Valid Slug
		slug, err := NewSlug("test-asset")
		assert.NoError(t, err)
		assert.Equal(t, "test-asset", slug.Value())

		// Invalid Slug
		_, err = NewSlug("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidSlug, err)

		_, err = NewSlug("invalid slug with spaces")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidSlug, err)
	})

	t.Run("Title", func(t *testing.T) {
		// Valid Title
		title, err := NewTitle("Test Asset Title")
		assert.NoError(t, err)
		assert.Equal(t, "Test Asset Title", title.Value())

		// Invalid Title
		_, err = NewTitle("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTitle, err)

		_, err = NewTitle("This is a very long title that exceeds the maximum allowed length of 100 characters and should cause an error")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTitle, err)
	})

	t.Run("Description", func(t *testing.T) {
		// Valid Description
		description, err := NewDescription("A test description")
		assert.NoError(t, err)
		assert.Equal(t, "A test description", description.Value())

		// Invalid Description
		_, err = NewDescription("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDescription, err)
	})

	t.Run("AssetType", func(t *testing.T) {
		// Valid AssetType
		assetType, err := NewAssetType("movie")
		assert.NoError(t, err)
		assert.Equal(t, "movie", assetType.Value())

		// Invalid AssetType
		_, err = NewAssetType("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAssetType, err)

		_, err = NewAssetType("invalid-type")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAssetType, err)
	})

	t.Run("Genre", func(t *testing.T) {
		// Valid Genre
		genre, err := NewGenre("action")
		assert.NoError(t, err)
		assert.Equal(t, "action", genre.Value())

		// Invalid Genre
		_, err = NewGenre("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidGenre, err)
	})

	t.Run("Genres", func(t *testing.T) {
		// Valid Genres
		genres, err := NewGenres([]string{"action", "drama", "thriller"})
		assert.NoError(t, err)
		assert.Len(t, genres.Values(), 3)

		// Empty Genres
		genres, err = NewGenres([]string{})
		assert.NoError(t, err)
		assert.Len(t, genres.Values(), 0)

		// Too Many Genres
		tooManyGenres := make([]string, 21)
		for i := range tooManyGenres {
			tooManyGenres[i] = "genre"
		}
		_, err = NewGenres(tooManyGenres)
		assert.Error(t, err)
		assert.Equal(t, ErrTooManyGenres, err)
	})

	t.Run("Tags", func(t *testing.T) {
		// Valid Tags
		tags, err := NewTags([]string{"tag1", "tag2", "tag3"})
		assert.NoError(t, err)
		assert.Len(t, tags.Values(), 3)

		// Empty Tags
		tags, err = NewTags([]string{})
		assert.NoError(t, err)
		assert.Len(t, tags.Values(), 0)

		// Too Many Tags
		tooManyTags := make([]string, 21)
		for i := range tooManyTags {
			tooManyTags[i] = "tag"
		}
		_, err = NewTags(tooManyTags)
		assert.Error(t, err)
		assert.Equal(t, ErrTooManyTags, err)
	})

	t.Run("OwnerID", func(t *testing.T) {
		// Valid OwnerID
		ownerID, err := NewOwnerID("user-123")
		assert.NoError(t, err)
		assert.Equal(t, "user-123", ownerID.Value())

		// Invalid OwnerID
		_, err = NewOwnerID("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidOwnerID, err)
	})

	t.Run("PublishRule", func(t *testing.T) {
		now := time.Now().UTC()
		future := now.Add(time.Hour)

		// Valid PublishRule
		publishRule, err := NewPublishRule(&now, &future, []string{"US", "CA"}, stringPtr("PG-13"))
		assert.NoError(t, err)
		assert.Equal(t, now, *publishRule.PublishAt())
		assert.Equal(t, future, *publishRule.UnpublishAt())
		assert.Equal(t, []string{"US", "CA"}, publishRule.Regions())
		assert.Equal(t, "PG-13", *publishRule.AgeRating())

		// Invalid PublishRule (publish after unpublish)
		_, err = NewPublishRule(&future, &now, []string{"US"}, nil)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPublishDates, err)

		// Too Many Regions
		tooManyRegions := make([]string, 51)
		for i := range tooManyRegions {
			tooManyRegions[i] = "US"
		}
		_, err = NewPublishRule(&now, nil, tooManyRegions, nil)
		assert.Error(t, err)
		assert.Equal(t, ErrTooManyRegions, err)
	})
}

func TestRichDomainModel(t *testing.T) {
	t.Run("CreateAssetWithValueObjects", func(t *testing.T) {
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)
		assert.Equal(t, "test-asset", asset.Slug().Value())
		assert.Equal(t, "Test Asset", asset.Title().Value())
		assert.Equal(t, "movie", asset.Type().Value())
		assert.Equal(t, constants.AssetStatusDraft, asset.Status())
	})

	t.Run("AssetLifecycleMethods", func(t *testing.T) {
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Test CanUpdateTitle
		assert.True(t, asset.CanUpdateTitle())

		// Test UpdateTitle
		newTitle, _ := NewTitle("Updated Title")
		err = asset.UpdateTitle(newTitle)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", asset.Title().Value())

		// Test CanUpdateDescription
		assert.True(t, asset.CanUpdateDescription())

		// Test UpdateDescription
		description, _ := NewDescription("Test description")
		err = asset.UpdateDescription(description)
		assert.NoError(t, err)
		assert.Equal(t, "Test description", asset.Description().Value())
	})

	t.Run("AssetPublishing", func(t *testing.T) {
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Test IsReadyForPublishing
		assert.False(t, asset.IsReadyForPublishing())

		// Add required content
		s3Object, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		video, err := asset.AddVideo("main", &VideoFormatMP4, *s3Object)
		assert.NoError(t, err)

		// Update video status to ready
		err = asset.UpdateVideoStatus(video.ID(), VideoStatusReady)
		assert.NoError(t, err)

		// Now should be ready for publishing
		assert.True(t, asset.IsReadyForPublishing())

		// Test CanBePublished
		assert.True(t, asset.CanBePublished())
	})

	t.Run("AssetHierarchy", func(t *testing.T) {
		parentSlug, _ := NewSlug("parent-asset")
		parentTitle, _ := NewTitle("Parent Asset")
		parentType, _ := NewAssetType("series")

		parent, err := NewAsset(*parentSlug, parentTitle, parentType)
		assert.NoError(t, err)

		childSlug, _ := NewSlug("child-asset")
		childTitle, _ := NewTitle("Child Asset")
		childType, _ := NewAssetType("episode")

		child, err := NewAsset(*childSlug, childTitle, childType)
		assert.NoError(t, err)

		// Set parent
		parentID := parent.ID()
		err = child.SetParentID(&parentID)
		assert.NoError(t, err)
		assert.Equal(t, parent.ID().Value(), child.ParentID().Value())
	})
}

func TestDomainServices(t *testing.T) {
	t.Run("AssetDomainService", func(t *testing.T) {
		service := NewDomainService(nil) // nil repo for testing

		// Test ValidateAssetForPublishing
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Should fail validation (no videos)
		err = service.ValidateAssetForPublishing(asset)
		assert.Error(t, err)

		// Add video and make it ready
		s3Object, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		video, err := asset.AddVideo("main", &VideoFormatMP4, *s3Object)
		assert.NoError(t, err)

		err = asset.UpdateVideoStatus(video.ID(), VideoStatusReady)
		assert.NoError(t, err)

		// Should pass validation
		err = service.ValidateAssetForPublishing(asset)
		assert.NoError(t, err)
	})

	t.Run("AssetPublishingService", func(t *testing.T) {
		domainService := NewDomainService(nil) // nil repo for testing
		service := NewPublishingService(domainService)

		// Test ValidatePublishingRules
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Add video and make it ready
		s3Object, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		video, err := asset.AddVideo("main", &VideoFormatMP4, *s3Object)
		assert.NoError(t, err)

		err = asset.UpdateVideoStatus(video.ID(), VideoStatusReady)
		assert.NoError(t, err)

		// Test with valid publish rule
		now := time.Now().UTC()
		publishRule, _ := NewPublishRule(&now, nil, []string{"US"}, nil)
		err = asset.SetPublishRule(publishRule)
		assert.NoError(t, err)

		err = service.ValidatePublishingRules(asset)
		assert.NoError(t, err)
	})

	t.Run("AssetMetrics", func(t *testing.T) {
		service := NewDomainService(nil) // nil repo for testing

		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Add videos
		s3Object1, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		_, err = asset.AddVideo("main", &VideoFormatMP4, *s3Object1)
		assert.NoError(t, err)

		s3Object2, _ := NewS3Object("test-bucket", "videos/trailer.m3u8", "https://test-bucket.s3.amazonaws.com/videos/trailer.m3u8")
		_, err = asset.AddVideo("trailer", &VideoFormatHLS, *s3Object2)
		assert.NoError(t, err)

		// Calculate metrics
		metrics := service.CalculateAssetMetrics(asset)
		assert.Equal(t, 2, metrics.TotalVideos)
		assert.Equal(t, 0, metrics.TotalImages)
		assert.Equal(t, 0, metrics.TotalCredits)
	})

	t.Run("StorageUsage", func(t *testing.T) {
		service := NewDomainService(nil) // nil repo for testing

		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Add videos with sizes
		s3Object1, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		video1, err := asset.AddVideo("main", &VideoFormatMP4, *s3Object1)
		assert.NoError(t, err)

		// Set video size
		video1.UpdateSize(1024 * 1024 * 100) // 100MB

		s3Object2, _ := NewS3Object("test-bucket", "videos/trailer.m3u8", "https://test-bucket.s3.amazonaws.com/videos/trailer.m3u8")
		video2, err := asset.AddVideo("trailer", &VideoFormatHLS, *s3Object2)
		assert.NoError(t, err)

		// Set video size
		video2.UpdateSize(1024 * 1024 * 50) // 50MB

		// Calculate storage usage
		usage := service.CalculateAssetStorageUsage(asset)
		assert.Equal(t, int64(1024*1024*150), usage.TotalStorage) // 150MB
		assert.Equal(t, int64(1024*1024*150), usage.VideoStorage) // 150MB
		assert.Equal(t, int64(0), usage.ImageStorage)             // 0MB
	})
}

func TestValueObjectEquality(t *testing.T) {
	t.Run("AssetIDEquality", func(t *testing.T) {
		id1, _ := NewAssetID("asset-123")
		id2, _ := NewAssetID("asset-123")
		id3, _ := NewAssetID("asset-456")

		assert.True(t, id1.Equals(*id2))
		assert.False(t, id1.Equals(*id3))
	})

	t.Run("SlugEquality", func(t *testing.T) {
		slug1, _ := NewSlug("test-asset")
		slug2, _ := NewSlug("test-asset")
		slug3, _ := NewSlug("different-asset")

		assert.True(t, slug1.Equals(*slug2))
		assert.False(t, slug1.Equals(*slug3))
	})

	t.Run("TitleEquality", func(t *testing.T) {
		title1, _ := NewTitle("Test Title")
		title2, _ := NewTitle("Test Title")
		title3, _ := NewTitle("Different Title")

		assert.True(t, title1.Equals(*title2))
		assert.False(t, title1.Equals(*title3))
	})
}

func TestComplexValueObjects(t *testing.T) {
	t.Run("S3Object", func(t *testing.T) {
		s3Object, err := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		assert.NoError(t, err)
		assert.Equal(t, "test-bucket", s3Object.Bucket())
		assert.Equal(t, "videos/main.mp4", s3Object.Key())
		assert.Equal(t, "https://test-bucket.s3.amazonaws.com/videos/main.mp4", s3Object.URL())

		// Invalid S3Object
		_, err = NewS3Object("", "key", "url")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidS3Bucket, err)
	})

	t.Run("StreamInfo", func(t *testing.T) {
		downloadURL := "https://example.com/download"
		cdnPrefix := "https://cdn.example.com"
		url := "https://example.com/stream"

		streamInfo, err := NewStreamInfo(&downloadURL, &cdnPrefix, &url)
		assert.NoError(t, err)
		assert.Equal(t, downloadURL, *streamInfo.DownloadURL())
		assert.Equal(t, cdnPrefix, *streamInfo.CDNPrefix())
		assert.Equal(t, url, *streamInfo.URL())
	})

	t.Run("TranscodingInfo", func(t *testing.T) {
		errorMsg := "Transcoding failed"
		completedAt := time.Now().UTC()

		transcodingInfo, err := NewTranscodingInfo("job-123", 75.5, "https://example.com/output.mp4", &errorMsg, &completedAt)
		assert.NoError(t, err)
		assert.Equal(t, "job-123", transcodingInfo.JobID())
		assert.Equal(t, 75.5, transcodingInfo.Progress())
		assert.Equal(t, "https://example.com/output.mp4", transcodingInfo.OutputURL())
		assert.Equal(t, errorMsg, *transcodingInfo.Error())
		assert.Equal(t, completedAt, *transcodingInfo.CompletedAt())
	})

	t.Run("Credit", func(t *testing.T) {
		personID := "person-123"
		credit, err := NewCredit("Director", "John Doe", &personID)
		assert.NoError(t, err)
		assert.Equal(t, "Director", credit.Role())
		assert.Equal(t, "John Doe", credit.Name())
		assert.Equal(t, personID, *credit.PersonID())

		// Invalid Credit
		_, err = NewCredit("", "John Doe", nil)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCreditRole, err)
	})
}

// Helper function for creating string pointers
func stringPtr(s string) *string {
	return &s
}
