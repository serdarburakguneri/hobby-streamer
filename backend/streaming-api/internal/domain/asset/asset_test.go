package asset

import (
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestAsset_IsPublished(t *testing.T) {
	now := time.Now().UTC()
	publishAt := now.Add(-time.Hour)
	unpublishAt := now.Add(time.Hour)
	pr, _ := NewPublishRuleValue(&publishAt, &unpublishAt, nil, nil)
	asset := &Asset{publishRule: pr}
	assert.True(t, asset.IsPublished())

	future := now.Add(time.Hour)
	pr2, _ := NewPublishRuleValue(&future, &unpublishAt, nil, nil)
	asset2 := &Asset{publishRule: pr2}
	assert.False(t, asset2.IsPublished())
}

func TestAsset_IsReady(t *testing.T) {
	video := Video{status: ptr(constants.VideoStatusReady)}
	asset := &Asset{videos: []Video{video}}
	assert.True(t, asset.IsReady())

	asset2 := &Asset{videos: []Video{}}
	assert.False(t, asset2.IsReady())
}

func TestAsset_GetMainVideo(t *testing.T) {
	mainType, _ := NewVideoType(constants.VideoTypeMain)
	video := Video{videoType: mainType}
	asset := &Asset{videos: []Video{video}}
	assert.NotNil(t, asset.GetMainVideo())

	asset2 := &Asset{videos: []Video{}}
	assert.Nil(t, asset2.GetMainVideo())
}

func TestAsset_GetThumbnail(t *testing.T) {
	thumbType, _ := NewImageType(constants.ImageTypeThumbnail)
	image := Image{imageType: thumbType}
	asset := &Asset{images: []Image{image}}
	assert.NotNil(t, asset.GetThumbnail())

	asset2 := &Asset{images: []Image{}}
	assert.Nil(t, asset2.GetThumbnail())
}

func TestAsset_IsAgeAppropriate(t *testing.T) {
	pr := &PublishRuleValue{ageRating: ptr(constants.AgeRatingPG13)}
	asset := &Asset{publishRule: pr}
	assert.True(t, asset.IsAgeAppropriate(15))
	assert.False(t, asset.IsAgeAppropriate(12))
}

func ptr[T any](v T) *T { return &v }
