package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

func (h *EventHandlers) HandleAnalyzeJobCompleted(ctx context.Context, ev *events.Event) error {
	var payload messages.JobCompletionPayload
	if err := unmarshalEventData(h.logger, ev, &payload); err != nil {
		return err
	}
	if payload.Success {
		if h.pipeline != nil {
			_ = h.pipeline.MarkCompleted(ctx, payload.AssetID, payload.VideoID, "analyze")
		}
		assetIDVO, err := valueobjects.NewAssetID(payload.AssetID)
		if err != nil {
			return err
		}
		cmd := commands.UpdateVideoMetadataCommand{
			AssetID:     *assetIDVO,
			VideoID:     payload.VideoID,
			Width:       payload.Width,
			Height:      payload.Height,
			Duration:    payload.Duration,
			Bitrate:     payload.Bitrate,
			Codec:       payload.Codec,
			Size:        payload.Size,
			ContentType: payload.ContentType,
		}
		if err := h.appService.UpdateVideoMetadata(ctx, cmd); err != nil {
			return err
		}
		// Do not auto-trigger HLS/DASH; user initiates via GraphQL mutation
	}
	return nil
}
