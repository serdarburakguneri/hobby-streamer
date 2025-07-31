package asset

import (
	"context"
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func TestValueObjects(t *testing.T) {
	t.Run("AssetID", func(t *testing.T) {
		assetID, err := valueobjects.NewAssetID("asset-123")
		assert.NoError(t, err)
		assert.Equal(t, "asset-123", assetID.Value())

		_, err = valueobjects.NewAssetID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "asset ID cannot be empty")
	})

	t.Run("Slug", func(t *testing.T) {
		slug, err := valueobjects.NewSlug("test-asset")
		assert.NoError(t, err)
		assert.Equal(t, "test-asset", slug.Value())

		_, err = valueobjects.NewSlug("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "slug cannot be empty")

		_, err = valueobjects.NewSlug("invalid slug with spaces")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "slug contains invalid characters")
	})

	t.Run("Title", func(t *testing.T) {
		title, err := valueobjects.NewTitle("Test Asset Title")
		assert.NoError(t, err)
		assert.Equal(t, "Test Asset Title", title.Value())

		_, err = valueobjects.NewTitle("")
		if err != nil {
			assert.Contains(t, err.Error(), "title cannot be empty")
		} else {
			t.Error("Expected error for empty title, got nil")
		}

		longTitle := "This is a very long title that exceeds the maximum allowed length of 200 characters and should cause an error. This is a very long title that exceeds the maximum allowed length of 200 characters and should cause an error. This is a very long title that exceeds the maximum allowed length of 200 characters and should cause an error."
		_, err = valueobjects.NewTitle(longTitle)
		if err != nil {
			assert.Contains(t, err.Error(), "title too long")
		} else {
			t.Error("Expected error for long title, got nil")
		}
	})

	t.Run("Description", func(t *testing.T) {
		description, err := valueobjects.NewDescription("A test description")
		assert.NoError(t, err)
		assert.Equal(t, "A test description", description.Value())

		_, err = valueobjects.NewDescription("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description cannot be empty after trimming")
	})

	t.Run("AssetType", func(t *testing.T) {
		assetType, err := valueobjects.NewAssetType("movie")
		assert.NoError(t, err)
		assert.Equal(t, "movie", assetType.Value())

		_, err = valueobjects.NewAssetType("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "asset type cannot be empty")

		_, err = valueobjects.NewAssetType("invalid-type")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid asset type")
	})

	t.Run("Genre", func(t *testing.T) {
		genre, err := valueobjects.NewGenre("action")
		assert.NoError(t, err)
		assert.Equal(t, "action", genre.Value())

		_, err = valueobjects.NewGenre("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "genre cannot be empty")

		_, err = valueobjects.NewGenre("invalid-genre")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid genre")
	})

	t.Run("Genres", func(t *testing.T) {
		genres, err := valueobjects.NewGenres([]string{"action", "drama", "thriller"})
		assert.NoError(t, err)
		assert.Len(t, genres.Values(), 3)

		genres, err = valueobjects.NewGenres([]string{})
		assert.NoError(t, err)
		assert.Len(t, genres.Values(), 0)

		tooManyGenres := make([]string, 21)
		for i := range tooManyGenres {
			tooManyGenres[i] = "genre"
		}
		_, err = valueobjects.NewGenres(tooManyGenres)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many genres")
	})

	t.Run("Tags", func(t *testing.T) {
		tags, err := valueobjects.NewTags([]string{"tag1", "tag2", "tag3"})
		assert.NoError(t, err)
		assert.Len(t, tags.Values(), 3)

		tags, err = valueobjects.NewTags([]string{})
		assert.NoError(t, err)
		assert.Len(t, tags.Values(), 0)

		tooManyTags := make([]string, 21)
		for i := range tooManyTags {
			tooManyTags[i] = "tag"
		}
		_, err = valueobjects.NewTags(tooManyTags)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many tags")
	})

	t.Run("OwnerID", func(t *testing.T) {
		ownerID, err := valueobjects.NewOwnerID("user-123")
		assert.NoError(t, err)
		assert.Equal(t, "user-123", ownerID.Value())

		_, err = valueobjects.NewOwnerID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "owner ID cannot be empty")
	})

	t.Run("PublishRule", func(t *testing.T) {
		now := time.Now().UTC()
		publishRule, err := valueobjects.NewPublishRule(&now, nil, []string{"US"}, nil)
		assert.NoError(t, err)
		assert.NotNil(t, publishRule)

		unpublishAt := now.Add(-time.Hour)
		_, err = valueobjects.NewPublishRule(&now, &unpublishAt, []string{"US"}, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publish date cannot be after unpublish date")
	})
}

func TestRichDomainModel(t *testing.T) {
	t.Run("CreateAssetWithValueObjects", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, err := entity.NewAsset(*slug, title, assetType)
		assert.NoError(t, err)
		assert.Equal(t, "test-asset", asset.Slug().Value())
		assert.Equal(t, "Test Asset", asset.Title().Value())
		assert.Equal(t, "movie", asset.Type().Value())
		assert.Equal(t, constants.AssetStatusDraft, asset.Status())
	})

	t.Run("AssetLifecycleMethods", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, err := entity.NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		assert.True(t, asset.CanUpdateTitle())

		newTitle, _ := valueobjects.NewTitle("Updated Title")
		asset.UpdateTitle(newTitle)
		assert.Equal(t, "Updated Title", asset.Title().Value())

		assert.True(t, asset.CanUpdateDescription())

		description, _ := valueobjects.NewDescription("Test description")
		asset.UpdateDescription(description)
		assert.Equal(t, "Test description", asset.Description().Value())
	})

	t.Run("AssetPublishing", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, err := entity.NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		assert.False(t, asset.IsReadyForPublishing())

		s3Object, _ := valueobjects.NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat := valueobjects.VideoFormat(constants.VideoStreamingFormatRaw)
		video, err := asset.AddVideo("main", &videoFormat, *s3Object)
		assert.NoError(t, err)

		err = asset.UpdateVideoStatus(video.ID().Value(), valueobjects.VideoStatus(constants.VideoStatusReady))
		assert.NoError(t, err)

		assert.True(t, asset.IsReadyForPublishing())

		assert.True(t, asset.CanBePublished())
	})

	t.Run("AssetHierarchy", func(t *testing.T) {
		parentSlug, _ := valueobjects.NewSlug("parent-asset")
		parentTitle, _ := valueobjects.NewTitle("Parent Asset")
		parentType, _ := valueobjects.NewAssetType("movie")

		parent, err := entity.NewAsset(*parentSlug, parentTitle, parentType)
		assert.NoError(t, err)

		childSlug, _ := valueobjects.NewSlug("child-asset")
		childTitle, _ := valueobjects.NewTitle("Child Asset")
		childType, _ := valueobjects.NewAssetType("movie")

		child, err := entity.NewAsset(*childSlug, childTitle, childType)
		assert.NoError(t, err)

		parentID := parent.ID()
		child.SetParentID(&parentID)
		assert.Equal(t, parent.ID().Value(), child.ParentID().Value())
	})
}

func TestDomainServices(t *testing.T) {
	t.Run("AssetDomainService", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, err := entity.NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		err = asset.ValidateForPublishing()
		assert.Error(t, err)

		s3Object, _ := valueobjects.NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat := valueobjects.VideoFormat(constants.VideoStreamingFormatRaw)
		video, err := asset.AddVideo("main", &videoFormat, *s3Object)
		assert.NoError(t, err)

		err = asset.UpdateVideoStatus(video.ID().Value(), valueobjects.VideoStatus(constants.VideoStatusReady))
		assert.NoError(t, err)

		now := time.Now().UTC()
		publishRule, _ := valueobjects.NewPublishRule(&now, nil, []string{"US"}, nil)
		err = asset.SetPublishRule(publishRule)
		assert.NoError(t, err)

		err = asset.ValidateForPublishing()
		assert.NoError(t, err)
	})

	t.Run("PublishingService", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, _ := entity.NewAsset(*slug, title, assetType)

		assert.False(t, asset.IsReadyForPublishing())
		assert.Error(t, asset.ValidateForPublishing())
	})

	t.Run("AssetMetrics", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, err := entity.NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		// Add videos
		s3Object1, _ := valueobjects.NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat1 := valueobjects.VideoFormat(constants.VideoStreamingFormatRaw)
		_, err = asset.AddVideo("main", &videoFormat1, *s3Object1)
		assert.NoError(t, err)

		s3Object2, _ := valueobjects.NewS3Object("test-bucket", "videos/trailer.m3u8", "https://test-bucket.s3.amazonaws.com/videos/trailer.m3u8")
		videoFormat2 := valueobjects.VideoFormat(constants.VideoStreamingFormatHLS)
		_, err = asset.AddVideo("trailer", &videoFormat2, *s3Object2)
		assert.NoError(t, err)

		metrics := asset.CalculateMetrics()
		metricsMap := metrics.(map[string]interface{})
		assert.Equal(t, 2, metricsMap["videoCount"])
		assert.Equal(t, 0, metricsMap["imageCount"])
		assert.Equal(t, 0, metricsMap["creditCount"])
	})

	t.Run("StorageUsage", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-asset")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, err := entity.NewAsset(*slug, title, assetType)
		assert.NoError(t, err)

		s3Object1, _ := valueobjects.NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		videoFormat1 := valueobjects.VideoFormat(constants.VideoStreamingFormatRaw)
		video1, err := asset.AddVideo("main", &videoFormat1, *s3Object1)
		assert.NoError(t, err)
		video1.UpdateSize(1024 * 1024 * 100)

		s3Object2, _ := valueobjects.NewS3Object("test-bucket", "videos/trailer.m3u8", "https://test-bucket.s3.amazonaws.com/videos/trailer.m3u8")
		videoFormat2 := valueobjects.VideoFormat(constants.VideoStreamingFormatHLS)
		video2, err := asset.AddVideo("trailer", &videoFormat2, *s3Object2)
		assert.NoError(t, err)
		video2.UpdateSize(1024 * 1024 * 50)

		usage := asset.CalculateStorageUsage()
		usageMap := usage.(map[string]interface{})
		assert.Equal(t, int64(1024*1024*150), usageMap["totalSize"])
		assert.Equal(t, 2, usageMap["videoCount"])
	})
}

func TestValueObjectEquality(t *testing.T) {
	t.Run("AssetIDEquality", func(t *testing.T) {
		id1, _ := valueobjects.NewAssetID("asset-123")
		id2, _ := valueobjects.NewAssetID("asset-123")
		id3, _ := valueobjects.NewAssetID("asset-456")

		assert.True(t, id1.Equals(*id2))
		assert.False(t, id1.Equals(*id3))
	})

	t.Run("SlugEquality", func(t *testing.T) {
		slug1, _ := valueobjects.NewSlug("test-asset")
		slug2, _ := valueobjects.NewSlug("test-asset")
		slug3, _ := valueobjects.NewSlug("different-asset")

		assert.True(t, slug1.Equals(*slug2))
		assert.False(t, slug1.Equals(*slug3))
	})

	t.Run("TitleEquality", func(t *testing.T) {
		title1, _ := valueobjects.NewTitle("Test Title")
		title2, _ := valueobjects.NewTitle("Test Title")
		title3, _ := valueobjects.NewTitle("Different Title")

		assert.True(t, title1.Equals(*title2))
		assert.False(t, title1.Equals(*title3))
	})
}

func TestComplexValueObjects(t *testing.T) {
	t.Run("S3Object", func(t *testing.T) {
		s3Object, err := valueobjects.NewS3Object("test-bucket", "videos/main.mp4", "https://test-bucket.s3.amazonaws.com/videos/main.mp4")
		assert.NoError(t, err)
		assert.Equal(t, "test-bucket", s3Object.Bucket())
		assert.Equal(t, "videos/main.mp4", s3Object.Key())
		assert.Equal(t, "https://test-bucket.s3.amazonaws.com/videos/main.mp4", s3Object.URL())

		_, err = valueobjects.NewS3Object("", "key", "url")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bucket cannot be empty")
	})

	t.Run("StreamInfo", func(t *testing.T) {
		downloadURL := "https://example.com/download"
		cdnPrefix := "https://cdn.example.com"
		url := "https://example.com/stream"

		streamInfo, err := valueobjects.NewStreamInfo(&downloadURL, &cdnPrefix, &url)
		assert.NoError(t, err)
		assert.Equal(t, downloadURL, *streamInfo.DownloadURL())
		assert.Equal(t, cdnPrefix, *streamInfo.CDNPrefix())
		assert.Equal(t, url, *streamInfo.URL())
	})

	t.Run("TranscodingInfo", func(t *testing.T) {
		contentType, _ := valueobjects.NewContentType("video/mp4")
		transcodingInfo := valueobjects.NewTranscodingInfo(1920, 1080, 120.5, 5000, "h264", 1024*1024*100, *contentType)
		assert.Equal(t, 1920, transcodingInfo.Width())
		assert.Equal(t, 1080, transcodingInfo.Height())
		assert.Equal(t, 120.5, transcodingInfo.Duration())
		assert.Equal(t, 5000, transcodingInfo.Bitrate())
		assert.Equal(t, "h264", transcodingInfo.Codec())
		assert.Equal(t, int64(1024*1024*100), transcodingInfo.Size())
		assert.Equal(t, "video/mp4", transcodingInfo.ContentType().Value())
	})

	t.Run("Credit", func(t *testing.T) {
		credit, err := valueobjects.NewCredit("Director", "John Doe", 1)
		assert.NoError(t, err)
		assert.Equal(t, "Director", credit.Role())
		assert.Equal(t, "John Doe", credit.Name())
		assert.Equal(t, 1, credit.Order())

		_, err = valueobjects.NewCredit("", "John Doe", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role cannot be empty")
	})
}

type mockRepo struct {
	findByIDFunc func(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error)
}

func (m *mockRepo) FindByID(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockRepo) Save(ctx context.Context, asset *entity.Asset) error { return nil }

func (m *mockRepo) FindBySlug(ctx context.Context, slug valueobjects.Slug) (*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) Update(ctx context.Context, asset *entity.Asset) error { return nil }

func (m *mockRepo) Delete(ctx context.Context, id valueobjects.AssetID) error { return nil }

func (m *mockRepo) List(ctx context.Context, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) Search(ctx context.Context, query string, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindByParentID(ctx context.Context, parentID valueobjects.AssetID, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindByType(ctx context.Context, assetType valueobjects.AssetType, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindByGenre(ctx context.Context, genre valueobjects.Genre, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindByTag(ctx context.Context, tag valueobjects.Tag, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindPublished(ctx context.Context, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) FindByPublishDate(ctx context.Context, from, to valueobjects.CreatedAt, limit *int, offset *int) ([]*entity.Asset, error) {
	return nil, nil
}

func (m *mockRepo) Count(ctx context.Context) (int64, error) { return 0, nil }

func (m *mockRepo) CountByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID) (int64, error) {
	return 0, nil
}

func (m *mockRepo) CountByType(ctx context.Context, assetType valueobjects.AssetType) (int64, error) {
	return 0, nil
}

func (m *mockRepo) Exists(ctx context.Context, id valueobjects.AssetID) (bool, error) {
	return false, nil
}

func (m *mockRepo) ExistsBySlug(ctx context.Context, slug valueobjects.Slug) (bool, error) {
	return false, nil
}

func TestValidateAssetHierarchy(t *testing.T) {
	domainServiceWithRepo := func(findByIDFunc func(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error)) DomainService {
		return NewDomainService(&mockRepo{findByIDFunc: findByIDFunc})
	}

	slug, _ := valueobjects.NewSlug("test-asset")
	title, _ := valueobjects.NewTitle("Test Asset")
	assetType, _ := valueobjects.NewAssetType("movie")
	asset, _ := entity.NewAsset(*slug, title, assetType)
	assetID := asset.ID()

	t.Run("parentID is nil", func(t *testing.T) {
		ds := domainServiceWithRepo(nil)
		ctx := context.Background()
		err := ds.ValidateAssetHierarchy(ctx, asset, nil)
		assert.NoError(t, err)
	})

	t.Run("asset is its own parent", func(t *testing.T) {
		ds := domainServiceWithRepo(nil)
		ctx := context.Background()
		err := ds.ValidateAssetHierarchy(ctx, asset, &assetID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "asset cannot be its own parent")
	})

	otherID, _ := valueobjects.NewAssetID("parent-asset")

	t.Run("parent asset not found (repo error)", func(t *testing.T) {
		ds := domainServiceWithRepo(func(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
			return nil, assert.AnError
		})
		ctx := context.Background()
		err := ds.ValidateAssetHierarchy(ctx, asset, otherID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent asset not found")
	})

	t.Run("parent asset not found (nil)", func(t *testing.T) {
		ds := domainServiceWithRepo(func(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
			return nil, nil
		})
		ctx := context.Background()
		err := ds.ValidateAssetHierarchy(ctx, asset, otherID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent asset not found")
	})

	t.Run("parent asset not published", func(t *testing.T) {
		parent := entity.ReconstructAsset(
			*otherID,
			*slug,
			title,
			nil,
			assetType,
			nil,
			nil,
			nil,
			asset.CreatedAt(),
			asset.UpdatedAt(),
			nil,
			nil,
			[]valueobjects.Image{},
			make(map[string]*entity.Video),
			[]valueobjects.Credit{},
			nil,
			map[string]interface{}{},
		)
		ds := domainServiceWithRepo(func(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
			return parent, nil
		})
		ctx := context.Background()
		err := ds.ValidateAssetHierarchy(ctx, asset, otherID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent asset is not published")
	})

	t.Run("parent asset is published", func(t *testing.T) {
		publishAt := time.Now().Add(-time.Hour)
		publishRule, _ := valueobjects.NewPublishRule(&publishAt, nil, []string{"US"}, nil)
		parent := entity.ReconstructAsset(
			*otherID,
			*slug,
			title,
			nil,
			assetType,
			nil,
			nil,
			nil,
			asset.CreatedAt(),
			asset.UpdatedAt(),
			nil,
			nil,
			[]valueobjects.Image{},
			make(map[string]*entity.Video),
			[]valueobjects.Credit{},
			publishRule,
			map[string]interface{}{},
		)
		ds := domainServiceWithRepo(func(ctx context.Context, id valueobjects.AssetID) (*entity.Asset, error) {
			return parent, nil
		})
		ctx := context.Background()
		err := ds.ValidateAssetHierarchy(ctx, asset, otherID)
		assert.NoError(t, err)
	})
}
