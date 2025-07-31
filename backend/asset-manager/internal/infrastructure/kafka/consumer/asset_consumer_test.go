package consumer

import (
	"context"
	"testing"

	appcommands "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAssetAppService struct {
	mock.Mock
}

func (m *MockAssetAppService) AddVideo(ctx context.Context, cmd appcommands.AddVideoCommand) (*assetentity.Video, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*assetentity.Video), args.Error(1)
}

func TestAssetEventConsumer_HandleJobCompletion(t *testing.T) {
	mockAssetAppService := new(MockAssetAppService)
	mockProducer := &events.Producer{}
	handlers := NewEventHandlers(mockAssetAppService, mockProducer, logger.Get())

	ctx := context.Background()

	t.Run("analyze job completion", func(t *testing.T) {
		event := &events.Event{
			Type: "analyze.job.completed",
			Data: []byte(`{
				"assetId": "asset-123",
				"videoId": "video-456",
				"success": true,
				"duration": 120.5,
				"width": 1920,
				"height": 1080,
				"bitrate": 5000000,
				"codec": "h264",
				"size": 1000000000,
				"contentType": "video/mp4"
			}`),
		}

		mockAssetAppService.On("UpdateVideoAnalysis", ctx, "asset-123", "video-456", mock.Anything).Return(nil)

		err := handlers.HandleAnalyzeJobCompleted(ctx, event)

		assert.NoError(t, err)
		mockAssetAppService.AssertExpectations(t)
	})

	t.Run("hls job completion", func(t *testing.T) {
		event := &events.Event{
			Type: "hls.job.completed",
			Data: []byte(`{
				"assetId": "asset-123",
				"videoId": "video-456",
				"success": true,
				"bucket": "content-east",
				"key": "asset-123/hls/main/playlist.m3u8",
				"url": "https://cdn.example.com/asset-123/hls/main/playlist.m3u8",
				"segmentCount": 10,
				"videoCodec": "h264",
				"audioCodec": "aac",
				"avgSegmentDuration": 6.0,
				"segments": ["segment1.ts", "segment2.ts"]
			}`),
		}

		mockAssetAppService.On("AddHLSVideo", ctx, "asset-123", "video-456", mock.AnythingOfType("map[string]interface {}")).Return(nil)

		err := handlers.HandleHLSJobCompleted(ctx, event)

		assert.NoError(t, err)
		mockAssetAppService.AssertExpectations(t)
	})
}
