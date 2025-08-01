package asset

import (
	"testing"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestAssetEvents(t *testing.T) {
	t.Run("AssetCreatedEvent", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-slug")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, _ := entity.NewAsset(*slug, title, assetType)

		event := events.NewAssetCreatedEvent(asset.ID().Value(), asset.Slug().Value(), asset.Title().Value(), asset.Type().Value(), "")
		cloudEvent := event.EventType()

		assert.Equal(t, "asset.created", cloudEvent)
		assert.Equal(t, asset.ID().Value(), event.AssetID())
		assert.Equal(t, asset.Slug().Value(), event.Slug())
		assert.Equal(t, asset.Title().Value(), event.Title())
		assert.Equal(t, asset.Type().Value(), event.AssetType())
	})

	t.Run("VideoAddedEvent", func(t *testing.T) {
		slug, _ := valueobjects.NewSlug("test-slug")
		title, _ := valueobjects.NewTitle("Test Asset")
		assetType, _ := valueobjects.NewAssetType("movie")

		asset, _ := entity.NewAsset(*slug, title, assetType)

		videoFormat := valueobjects.VideoFormat("hls")
		s3Object, _ := valueobjects.NewS3Object("test-bucket", "test-key", "https://test-bucket.s3.amazonaws.com/test-key")
		video, _ := asset.AddVideo("main", &videoFormat, *s3Object, 1920, 1080, 120.5, 5000000, "h264", 1024000000, "video/mp4", "h264", "aac", "30fps", 2, 48000, nil)

		event := events.NewVideoAddedEvent(asset.ID().Value(), video.ID().Value(), video.Label().Value(), string(video.Format()))
		cloudEvent := event.EventType()

		assert.Equal(t, "video.added", cloudEvent)
		assert.Equal(t, asset.ID().Value(), event.AssetID())
		assert.Equal(t, video.ID().Value(), event.VideoID())
		assert.Equal(t, video.Label().Value(), event.Label())
		assert.Equal(t, string(video.Format()), event.Format())
	})
}
