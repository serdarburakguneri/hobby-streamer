package consumer

import (
	"context"
	"path"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

func (h *EventHandlers) HandleTranscodeHlsJobCompleted(ctx context.Context, ev *events.Event) error {
	var payload messages.JobCompletionPayload
	if err := unmarshalEventData(h.logger, ev, &payload); err != nil {
		return err
	}
	if !payload.Success {
		assetIDVO, err := valueobjects.NewAssetID(payload.AssetID)
		if err != nil {
			return err
		}
		statusFailed := valueobjects.VideoStatusFailed
		formatVO, _ := valueobjects.NewVideoFormat(string(valueobjects.VideoFormatHLS))
		h.appService.UpsertVideo(ctx, commands.UpsertVideoCommand{
			AssetID:       *assetIDVO,
			Label:         path.Base(payload.Key),
			Format:        formatVO,
			ContentType:   payload.ContentType,
			InitialStatus: &statusFailed,
		})
		if h.pipeline != nil {
			h.pipeline.MarkFailed(ctx, payload.AssetID, payload.VideoID, "hls", payload.Error)
		}

		ev2 := events.NewVideoStatusUpdatedEvent(payload.AssetID, payload.VideoID, statusFailed.Value())
		ev2.SetSource("asset-manager").SetEventVersion("1").SetCorrelationID(ev.CorrelationID).SetCausationID(ev.ID)
		h.publisher.Publish(ctx, events.AssetEventsTopic, ev2)
		return nil
	}
	assetIDVO, err := valueobjects.NewAssetID(payload.AssetID)
	if err != nil {
		return err
	}
	s3Obj, err := valueobjects.NewS3ObjectFromURL(payload.URL)
	if err != nil {
		return err
	}
	formatVO, err := valueobjects.NewVideoFormat(string(valueobjects.VideoFormatHLS))
	if err != nil {
		return err
	}
	cdnPrefix, playURL := h.cdn.BuildPlayURL(payload.Key)
	si, _ := valueobjects.NewStreamInfo(nil, &cdnPrefix, &playURL)

	statusReady := valueobjects.VideoStatusReady
	_, _, err = h.appService.UpsertVideo(ctx, commands.UpsertVideoCommand{
		AssetID:            *assetIDVO,
		Label:              path.Base(payload.Key),
		Format:             formatVO,
		StorageLocation:    *s3Obj,
		ContentType:        payload.ContentType,
		Codec:              payload.VideoCodec,
		VideoCodec:         payload.VideoCodec,
		AudioCodec:         payload.AudioCodec,
		FrameRate:          payload.FrameRate,
		AudioChannels:      payload.AudioChannels,
		AudioSampleRate:    payload.AudioSampleRate,
		Duration:           payload.Duration,
		Bitrate:            payload.Bitrate,
		Width:              payload.Width,
		Height:             payload.Height,
		Size:               payload.Size,
		StreamInfo:         si,
		SegmentCount:       payload.SegmentCount,
		AvgSegmentDuration: payload.AvgSegmentDuration,
		Segments:           payload.Segments,
		InitialStatus:      &statusReady,
	})
	if err != nil {
		return err
	}
	if h.pipeline != nil {
		h.pipeline.MarkCompleted(ctx, payload.AssetID, payload.VideoID, "hls")
	}

	ev2 := events.NewVideoStatusUpdatedEvent(payload.AssetID, payload.VideoID, valueobjects.VideoStatusReady.Value())
	ev2.SetSource("asset-manager").SetEventVersion("1").SetCorrelationID(ev.CorrelationID).SetCausationID(ev.ID)
	h.publisher.Publish(ctx, events.AssetEventsTopic, ev2)
	return nil
}
