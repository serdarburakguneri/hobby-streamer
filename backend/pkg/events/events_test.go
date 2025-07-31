package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent("test.event", map[string]interface{}{
		"key": "value",
	})

	assert.NotEmpty(t, event.ID)
	assert.Equal(t, "test.event", event.Type)
	assert.Equal(t, CloudEventsVersion, event.SpecVersion)
	assert.Equal(t, ContentTypeJSON, event.DataContentType)
	assert.False(t, event.Time.IsZero())
	assert.Equal(t, "value", event.Data.(map[string]interface{})["key"])
}

func TestEventValidation(t *testing.T) {
	event := NewEvent("test.event", nil)
	event.SetSource("test-service")
	assert.NoError(t, event.Validate())

	event.ID = ""
	assert.Error(t, event.Validate())

	event = NewEvent("test.event", nil)
	event.SetSource("test-service")
	event.Type = ""
	assert.Error(t, event.Validate())

	event = NewEvent("test.event", nil)
	event.SetSource("test-service")
	event.SpecVersion = ""
	assert.Error(t, event.Validate())

	event = NewEvent("test.event", nil)
	event.SetSource("test-service")
	event.Time = time.Time{}
	assert.Error(t, event.Validate())

	event = NewEvent("test.event", nil)
	event.Source = ""
	assert.Error(t, event.Validate())
}

func TestEventJSONMarshaling(t *testing.T) {
	event := NewEvent("test.event", map[string]interface{}{
		"key": "value",
	})
	event.SetSource("test-service")
	event.SetCorrelationID("corr-123")
	event.SetCausationID("cause-456")
	event.AddExtension("custom", "extension")

	jsonBytes, err := json.Marshal(event)
	require.NoError(t, err)

	var unmarshaledEvent Event
	err = json.Unmarshal(jsonBytes, &unmarshaledEvent)
	require.NoError(t, err)

	assert.Equal(t, event.ID, unmarshaledEvent.ID)
	assert.Equal(t, event.Type, unmarshaledEvent.Type)
	assert.Equal(t, event.Source, unmarshaledEvent.Source)
	assert.Equal(t, event.CorrelationID, unmarshaledEvent.CorrelationID)
	assert.Equal(t, event.CausationID, unmarshaledEvent.CausationID)
	assert.Equal(t, "extension", unmarshaledEvent.Extensions["custom"])
}

func TestEventHelperFunctions(t *testing.T) {
	assetEvent := NewAssetCreatedEvent("asset-123", "test-asset", "Test Asset", "movie")
	assert.Equal(t, AssetCreatedEventType, assetEvent.Type)
	assert.Equal(t, "asset-123", assetEvent.Data.(map[string]interface{})["assetId"])

	videoEvent := NewVideoAddedEvent("asset-123", "video-456", "Main Video", "hls")
	assert.Equal(t, VideoAddedEventType, videoEvent.Type)
	assert.Equal(t, "video-456", videoEvent.Data.(map[string]interface{})["videoId"])

	bucketEvent := NewBucketCreatedEvent("bucket-123", "Test Bucket", "test-bucket")
	assert.Equal(t, BucketCreatedEventType, bucketEvent.Type)
	assert.Equal(t, "bucket-123", bucketEvent.Data.(map[string]interface{})["bucketId"])

	jobEvent := NewJobAnalyzeRequestedEvent("asset-123", "video-456", "s3://bucket/video.mp4")
	assert.Equal(t, JobAnalyzeRequestedEventType, jobEvent.Type)
	assert.Equal(t, "analyze", jobEvent.Data.(map[string]interface{})["jobType"])
}

func TestGetDataAs(t *testing.T) {
	event := NewEvent("test.event", map[string]interface{}{
		"assetId": "asset-123",
		"title":   "Test Asset",
	})

	var data map[string]interface{}
	err := event.GetDataAs(&data)
	require.NoError(t, err)

	assert.Equal(t, "asset-123", data["assetId"])
	assert.Equal(t, "Test Asset", data["title"])
}

func TestEventString(t *testing.T) {
	event := NewEvent("test.event", map[string]interface{}{
		"key": "value",
	})

	str := event.String()
	assert.Contains(t, str, event.ID)
	assert.Contains(t, str, event.Type)
}
