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
		// Pre-create placeholders for HLS and DASH as transcoding
		statusTranscoding := valueobjects.VideoStatusTranscoding
		input := payload.URL
		hlsKey := path.Join(payload.AssetID, payload.VideoID, "hls", path.Base(payload.Key))
		hlsS3, _ := valueobjects.NewS3Object(payload.Bucket, hlsKey, "")
		hlsFmt, _ := valueobjects.NewVideoFormat(string(valueobjects.VideoFormatHLS))
		_, _, _ = h.appService.UpsertVideo(ctx, commands.UpsertVideoCommand{
			AssetID:         *assetIDVO,
			Label:           path.Base(hlsKey),
			Format:          hlsFmt,
			StorageLocation: *hlsS3,
			Duration:        0,
			Bitrate:         0,
			Width:           0,
			Height:          0,
			Size:            0,
			ContentType:     "application/x-mpegURL",
			InitialStatus:   &statusTranscoding,
		})
		input = payload.URL
		hlsEvt := events.NewJobTranscodeRequestedEvent(payload.AssetID, payload.VideoID, input, valueobjects.VideoFormatHLS.Value(), payload.Bucket, hlsKey)
		hlsEvt.SetSource("asset-manager").SetEventVersion("1").SetCorrelationID(events.BuildJobCorrelationID(payload.AssetID, payload.VideoID, "transcode", valueobjects.VideoFormatHLS.Value(), "main")).SetCausationID(ev.ID)
		if err := h.producer.SendEvent(ctx, events.HLSJobRequestedTopic, hlsEvt); err != nil {
			return err
		}
		dashKey := path.Join(payload.AssetID, payload.VideoID, "dash", path.Base(payload.Key))
		dashS3, _ := valueobjects.NewS3Object(payload.Bucket, dashKey, "")
		dashFmt, _ := valueobjects.NewVideoFormat(string(valueobjects.VideoFormatDASH))
		_, _, _ = h.appService.UpsertVideo(ctx, commands.UpsertVideoCommand{
			AssetID:         *assetIDVO,
			Label:           path.Base(dashKey),
			Format:          dashFmt,
			StorageLocation: *dashS3,
			ContentType:     "application/dash+xml",
			InitialStatus:   &statusTranscoding,
		})
		dashEvt := events.NewJobTranscodeRequestedEvent(payload.AssetID, payload.VideoID, input, valueobjects.VideoFormatDASH.Value(), payload.Bucket, dashKey)
		dashEvt.SetSource("asset-manager").SetEventVersion("1").SetCorrelationID(events.BuildJobCorrelationID(payload.AssetID, payload.VideoID, "transcode", valueobjects.VideoFormatDASH.Value(), "main")).SetCausationID(ev.ID)
		if err := h.producer.SendEvent(ctx, events.DASHJobRequestedTopic, dashEvt); err != nil {
			return err
		}
	}
	return nil
}
