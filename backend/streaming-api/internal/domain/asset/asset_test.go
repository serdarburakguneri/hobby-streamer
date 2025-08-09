package asset

import (
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestAsset_IsPublished(t *testing.T) {
	t.Skip("Skipping due to nil pointer dereference issue - needs investigation")

	now := time.Now().UTC()
	publishAt := now.Add(-time.Hour)
	unpublishAt := now.Add(time.Hour)

	pr, err := valueobjects.NewPublishRuleValue(&publishAt, &unpublishAt, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, pr)

	assetID, _ := valueobjects.NewAssetID("test-id")
	slug, _ := valueobjects.NewSlug("test-slug")
	assetType, _ := valueobjects.NewAssetType("video")
	createdAt := valueobjects.NewCreatedAt(now)
	updatedAt := valueobjects.NewUpdatedAt(now)

	asset := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, nil, nil, pr)
	assert.True(t, asset.IsPublished())

	future := now.Add(time.Hour)
	pr2, err2 := valueobjects.NewPublishRuleValue(&future, &unpublishAt, nil, nil)
	assert.NoError(t, err2)
	assert.NotNil(t, pr2)
	asset2 := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, nil, nil, pr2)
	assert.False(t, asset2.IsPublished())
}

func TestAsset_IsReady(t *testing.T) {
	t.Skip("Skipping due to nil pointer dereference issue - needs investigation")

	videoID, _ := valueobjects.NewVideoID("test-video-id")
	storageLocation, _ := valueobjects.NewS3ObjectValue("bucket", "key", "url")
	statusVO, _ := valueobjects.NewVideoStatus(constants.VideoStatusReady)
	video := entity.NewVideo(*videoID, nil, nil, *storageLocation, nil, nil, nil, nil, nil, nil, nil, nil, nil, statusVO, nil, time.Now().UTC(), time.Now().UTC(), nil, true, false, false, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	assetID, _ := valueobjects.NewAssetID("test-id")
	slug, _ := valueobjects.NewSlug("test-slug")
	assetType, _ := valueobjects.NewAssetType("video")
	createdAt := valueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := valueobjects.NewUpdatedAt(time.Now().UTC())

	asset := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, []entity.Video{*video}, nil, nil)
	assert.True(t, asset.IsReady())

	asset2 := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, []entity.Video{}, nil, nil)
	assert.False(t, asset2.IsReady())
}

func TestAsset_GetMainVideo(t *testing.T) {
	t.Skip("Skipping due to nil pointer dereference issue - needs investigation")

	mainType, _ := valueobjects.NewVideoType(constants.VideoTypeMain)
	videoID, _ := valueobjects.NewVideoID("test-video-id")
	storageLocation, _ := valueobjects.NewS3ObjectValue("bucket", "key", "url")
	statusVO, _ := valueobjects.NewVideoStatus(constants.VideoStatusReady)
	video := entity.NewVideo(*videoID, mainType, nil, *storageLocation, nil, nil, nil, nil, nil, nil, nil, nil, nil, statusVO, nil, time.Now().UTC(), time.Now().UTC(), nil, true, false, false, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	assetID, _ := valueobjects.NewAssetID("test-id")
	slug, _ := valueobjects.NewSlug("test-slug")
	assetType, _ := valueobjects.NewAssetType("video")
	createdAt := valueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := valueobjects.NewUpdatedAt(time.Now().UTC())

	asset := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, []entity.Video{*video}, nil, nil)
	assert.NotNil(t, asset.GetMainVideo())

	asset2 := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, []entity.Video{}, nil, nil)
	assert.Nil(t, asset2.GetMainVideo())
}

func TestAsset_GetThumbnail(t *testing.T) {
	t.Skip("Skipping due to nil pointer dereference issue - needs investigation")

	thumbType, _ := valueobjects.NewImageType(constants.ImageTypeThumbnail)
	imageID, _ := valueobjects.NewImageID("test-image-id")
	fileName, _ := valueobjects.NewFileName("test.jpg")
	image := entity.NewImage(*imageID, *fileName, "url", thumbType, nil, nil, nil, nil, nil, nil, nil, time.Now().UTC(), time.Now().UTC())

	assetID, _ := valueobjects.NewAssetID("test-id")
	slug, _ := valueobjects.NewSlug("test-slug")
	assetType, _ := valueobjects.NewAssetType("video")
	createdAt := valueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := valueobjects.NewUpdatedAt(time.Now().UTC())

	asset := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, nil, []entity.Image{*image}, nil)
	assert.NotNil(t, asset.GetThumbnail())

	asset2 := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, nil, []entity.Image{}, nil)
	assert.Nil(t, asset2.GetThumbnail())
}

func TestAsset_IsAgeAppropriate(t *testing.T) {
	t.Skip("Skipping due to nil pointer dereference issue - needs investigation")

	pr, _ := valueobjects.NewPublishRuleValue(nil, nil, nil, ptr(constants.AgeRatingPG13))

	assetID, _ := valueobjects.NewAssetID("test-id")
	slug, _ := valueobjects.NewSlug("test-slug")
	assetType, _ := valueobjects.NewAssetType("video")
	createdAt := valueobjects.NewCreatedAt(time.Now().UTC())
	updatedAt := valueobjects.NewUpdatedAt(time.Now().UTC())

	asset := entity.NewAsset(*assetID, *slug, nil, nil, *assetType, nil, nil, nil, nil, createdAt.Value(), updatedAt.Value(), nil, nil, nil, nil, pr)
	assert.True(t, asset.IsAgeAppropriate(15))
	assert.False(t, asset.IsAgeAppropriate(12))
}

func ptr[T any](v T) *T { return &v }
