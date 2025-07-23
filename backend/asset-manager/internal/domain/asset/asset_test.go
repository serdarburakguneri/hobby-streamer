package asset

import (
	"context"
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func TestValueObjects(t *testing.T) {
	t.Run("AssetID", func(t *testing.T) {
		// Valid AssetID
		assetID, err := NewAssetID("asset-123")
		assert.NoError(t, err)
		assert.Equal(t, "asset-123", assetID.Value())

		// Invalid AssetID
		_, err = NewAssetID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid asset ID")
	})

	t.Run("Slug", func(t *testing.T) {
		// Valid Slug
		slug, err := NewSlug("test-asset")
		assert.NoError(t, err)
		assert.Equal(t, "test-asset", slug.Value())

		// Invalid Slug
		_, err = NewSlug("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid slug")

		_, err = NewSlug("invalid slug with spaces")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid slug")
	})

	t.Run("Title", func(t *testing.T) {
		// Valid Title
		title, err := NewTitle("Test Asset Title")
		assert.NoError(t, err)
		assert.Equal(t, "Test Asset Title", title.Value())

		// Invalid Title
		_, err = NewTitle("")
		if err != nil {
			assert.Contains(t, err.Error(), "invalid title")
		} else {
			t.Error("Expected error for empty title, got nil")
		}

		longTitle := "This is a very long title that exceeds the maximum allowed length of 200 characters and should cause an error. This is a very long title that exceeds the maximum allowed length of 200 characters and should cause an error. This is a very long title that exceeds the maximum allowed length of 200 characters and should cause an error."
		_, err = NewTitle(longTitle)
		if err != nil {
			assert.Contains(t, err.Error(), "invalid title")
		} else {
			t.Error("Expected error for long title, got nil")
		}
	})

	t.Run("Description", func(t *testing.T) {
		// Valid Description
		description, err := NewDescription("A test description")
		assert.NoError(t, err)
		assert.Equal(t, "A test description", description.Value())

		// Invalid Description
		_, err = NewDescription("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid description")
	})

	t.Run("AssetType", func(t *testing.T) {
		// Valid AssetType
		assetType, err := NewAssetType("movie")
		assert.NoError(t, err)
		assert.Equal(t, "movie", assetType.Value())

		// Invalid AssetType
		_, err = NewAssetType("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid asset type")

		_, err = NewAssetType("invalid-type")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid asset type")
	})

	t.Run("Genre", func(t *testing.T) {
		// Valid Genre
		genre, err := NewGenre("action")
		assert.NoError(t, err)
		assert.Equal(t, "action", genre.Value())

		// Invalid Genre
		_, err = NewGenre("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid genre")
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
		assert.Contains(t, err.Error(), "too many genres")
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
		assert.Contains(t, err.Error(), "too many tags")
	})

	t.Run("OwnerID", func(t *testing.T) {
		// Valid OwnerID
		ownerID, err := NewOwnerID("user-123")
		assert.NoError(t, err)
		assert.Equal(t, "user-123", ownerID.Value())

		// Invalid OwnerID
		_, err = NewOwnerID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid owner ID")
	})

	t.Run("PublishRule", func(t *testing.T) {
		// Valid PublishRule
		now := time.Now().UTC()
		publishRule, err := NewPublishRule(&now, nil, []string{"US"}, nil)
		assert.NoError(t, err)
		assert.NotNil(t, publishRule)

		// Invalid PublishRule (publish date after unpublish date)
		unpublishAt := now.Add(-time.Hour)
		_, err = NewPublishRule(&now, &unpublishAt, []string{"US"}, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid publish dates")
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
		videoFormat := VideoFormat(constants.VideoStreamingFormatRaw)
		video, err := asset.AddVideo("main", &videoFormat, *s3Object)
		assert.NoError(t, err)

		// Update video status to ready
		err = asset.UpdateVideoStatus(video.ID(), VideoStatus(constants.VideoStatusReady))
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
		// Test ValidateAssetForPublishing
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Should fail validation (no videos)
		err = asset.ValidateForPublishing()
		assert.Error(t, err)

		// Add video and make it ready
		s3Object, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat := VideoFormat(constants.VideoStreamingFormatRaw)
		video, err := asset.AddVideo("main", &videoFormat, *s3Object)
		assert.NoError(t, err)

		err = asset.UpdateVideoStatus(video.ID(), VideoStatus(constants.VideoStatusReady))
		assert.NoError(t, err)

		// Add publish rule
		now := time.Now().UTC()
		publishRule, _ := NewPublishRule(&now, nil, []string{"US"}, nil)
		err = asset.SetPublishRule(publishRule)
		assert.NoError(t, err)

		// Should pass validation
		err = asset.ValidateForPublishing()
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
		videoFormat := VideoFormat(constants.VideoStreamingFormatRaw)
		video, err := asset.AddVideo("main", &videoFormat, *s3Object)
		assert.NoError(t, err)

		err = asset.UpdateVideoStatus(video.ID(), VideoStatus(constants.VideoStatusReady))
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
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Add videos
		s3Object1, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat1 := VideoFormat(constants.VideoStreamingFormatRaw)
		_, err = asset.AddVideo("main", &videoFormat1, *s3Object1)
		assert.NoError(t, err)

		s3Object2, _ := NewS3Object("test-bucket", "videos/trailer.m3u8", "https://test-bucket.s3.amazonaws.com/videos/trailer.m3u8")
		videoFormat2 := VideoFormat(constants.VideoStreamingFormatHLS)
		_, err = asset.AddVideo("trailer", &videoFormat2, *s3Object2)
		assert.NoError(t, err)

		// Calculate metrics
		metrics := asset.CalculateMetrics()
		assert.Equal(t, 2, metrics.TotalVideos)
		assert.Equal(t, 0, metrics.TotalImages)
		assert.Equal(t, 0, metrics.TotalCredits)
	})

	t.Run("StorageUsage", func(t *testing.T) {
		slug, _ := NewSlug("test-asset")
		title, _ := NewTitle("Test Asset")
		assetType, _ := NewAssetType("movie")

		asset, err := NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Add videos with sizes
		s3Object1, _ := NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat1 := VideoFormat(constants.VideoStreamingFormatRaw)
		video1, err := asset.AddVideo("main", &videoFormat1, *s3Object1)
		assert.NoError(t, err)

		// Set video size
		video1.UpdateSize(1024 * 1024 * 100) // 100MB

		s3Object2, _ := NewS3Object("test-bucket", "videos/trailer.m3u8", "https://test-bucket.s3.amazonaws.com/videos/trailer.m3u8")
		videoFormat2 := VideoFormat(constants.VideoStreamingFormatHLS)
		video2, err := asset.AddVideo("trailer", &videoFormat2, *s3Object2)
		assert.NoError(t, err)

		// Set video size
		video2.UpdateSize(1024 * 1024 * 50) // 50MB

		// Calculate storage usage
		usage := asset.CalculateStorageUsage()
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
		assert.Contains(t, err.Error(), "invalid S3 bucket")
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
		assert.Contains(t, err.Error(), "invalid credit role")
	})
}

type mockRepo struct {
	findByIDFunc func(ctx context.Context, id AssetID) (*Asset, error)
}

func (m *mockRepo) FindByID(ctx context.Context, id AssetID) (*Asset, error) {
	return m.findByIDFunc(ctx, id)
}
func (m *mockRepo) FindByIDs(ctx context.Context, ids []AssetID) ([]*Asset, error) { return nil, nil }
func (m *mockRepo) Create(ctx context.Context, asset *Asset) error                 { return nil }
func (m *mockRepo) Update(ctx context.Context, asset *Asset) error                 { return nil }
func (m *mockRepo) Delete(ctx context.Context, id AssetID) error                   { return nil }
func (m *mockRepo) List(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	return nil, nil
}
func (m *mockRepo) Search(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	return nil, nil
}
func (m *mockRepo) GetByOwnerID(ctx context.Context, ownerID string, limit *int, lastKey map[string]interface{}) ([]*Asset, error) {
	return nil, nil
}
func (m *mockRepo) AddChild(ctx context.Context, parentID, childID AssetID) error    { return nil }
func (m *mockRepo) RemoveChild(ctx context.Context, parentID, childID AssetID) error { return nil }
func (m *mockRepo) GetChildren(ctx context.Context, parentID AssetID) ([]*Asset, error) {
	return nil, nil
}
func (m *mockRepo) HasChild(ctx context.Context, parentID, childID AssetID) (bool, error) {
	return false, nil
}
func (m *mockRepo) FindBySlug(ctx context.Context, slug Slug) (*Asset, error) { return nil, nil }
func (m *mockRepo) FindByTypeAndGenre(ctx context.Context, assetType *AssetType, genre *Genre) ([]*Asset, error) {
	return nil, nil
}
func (m *mockRepo) FindChildren(ctx context.Context, parentID AssetID) ([]*Asset, error) {
	return nil, nil
}
func (m *mockRepo) FindParent(ctx context.Context, childID AssetID) (*Asset, error) { return nil, nil }
func (m *mockRepo) Save(ctx context.Context, asset *Asset) error                    { return nil }

func TestValidateAssetHierarchy(t *testing.T) {
	domainServiceWithRepo := func(findByIDFunc func(ctx context.Context, id AssetID) (*Asset, error)) *DomainService {
		return NewDomainService(&mockRepo{findByIDFunc: findByIDFunc})
	}

	slug, _ := NewSlug("test-asset")
	title, _ := NewTitle("Test Asset")
	assetType, _ := NewAssetType("movie")
	asset, _ := NewAsset(*slug, title, assetType)
	assetID := asset.ID()

	t.Run("parentID is nil", func(t *testing.T) {
		ds := domainServiceWithRepo(nil)
		err := ds.ValidateAssetHierarchy(asset, nil)
		assert.NoError(t, err)
	})

	t.Run("asset is its own parent", func(t *testing.T) {
		ds := domainServiceWithRepo(nil)
		err := ds.ValidateAssetHierarchy(asset, &assetID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "asset cannot be its own parent")
	})

	otherID, _ := NewAssetID("parent-asset")

	t.Run("parent asset not found (repo error)", func(t *testing.T) {
		ds := domainServiceWithRepo(func(ctx context.Context, id AssetID) (*Asset, error) {
			return nil, assert.AnError
		})
		err := ds.ValidateAssetHierarchy(asset, otherID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent asset not found")
	})

	t.Run("parent asset not found (nil)", func(t *testing.T) {
		ds := domainServiceWithRepo(func(ctx context.Context, id AssetID) (*Asset, error) {
			return nil, nil
		})
		err := ds.ValidateAssetHierarchy(asset, otherID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent asset not found")
	})

	t.Run("parent asset not published", func(t *testing.T) {
		parent := ReconstructAsset(
			*otherID,
			*slug,
			title,
			nil, // description
			assetType,
			nil, // genre
			nil, // genres
			nil, // tags
			asset.CreatedAt(),
			asset.UpdatedAt(),
			nil,       // ownerID
			nil,       // parentID
			[]Image{}, // images
			make(map[string]*Video),
			[]Credit{}, // credits
			nil,        // publishRule
			map[string]interface{}{},
		)
		// Simulate draft status by leaving publishRule nil
		ds := domainServiceWithRepo(func(ctx context.Context, id AssetID) (*Asset, error) {
			return parent, nil
		})
		err := ds.ValidateAssetHierarchy(asset, otherID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent asset must be published")
	})

	t.Run("parent asset is published", func(t *testing.T) {
		publishAt := time.Now().Add(-time.Hour)
		publishRule, _ := NewPublishRule(&publishAt, nil, []string{"US"}, nil)
		parent := ReconstructAsset(
			*otherID,
			*slug,
			title,
			nil, // description
			assetType,
			nil, // genre
			nil, // genres
			nil, // tags
			asset.CreatedAt(),
			asset.UpdatedAt(),
			nil,       // ownerID
			nil,       // parentID
			[]Image{}, // images
			make(map[string]*Video),
			[]Credit{}, // credits
			publishRule,
			map[string]interface{}{},
		)
		ds := domainServiceWithRepo(func(ctx context.Context, id AssetID) (*Asset, error) {
			return parent, nil
		})
		err := ds.ValidateAssetHierarchy(asset, otherID)
		assert.NoError(t, err)
	})
}
