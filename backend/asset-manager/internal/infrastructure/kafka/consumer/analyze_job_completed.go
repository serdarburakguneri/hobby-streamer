package consumer

import (
	"context"
	"path"

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
		input := payload.URL
		hlsKey := path.Join(payload.AssetID, payload.VideoID, "hls", path.Base(payload.Key))
		hlsEvt := events.NewJobTranscodeRequestedEvent(payload.AssetID, payload.VideoID, input, valueobjects.VideoFormatHLS.Value(), payload.Bucket, hlsKey)
		if err := h.producer.SendEvent(ctx, events.HLSJobRequestedTopic, hlsEvt); err != nil {
			return err
		}
		dashKey := path.Join(payload.AssetID, payload.VideoID, "dash", path.Base(payload.Key))
		dashEvt := events.NewJobTranscodeRequestedEvent(payload.AssetID, payload.VideoID, input, valueobjects.VideoFormatDASH.Value(), payload.Bucket, dashKey)
		if err := h.producer.SendEvent(ctx, events.DASHJobRequestedTopic, dashEvt); err != nil {
			return err
		}
	}
	return nil
}
