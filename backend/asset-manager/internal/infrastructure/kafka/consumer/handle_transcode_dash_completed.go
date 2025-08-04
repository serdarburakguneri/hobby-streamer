package consumer

import (
	"context"
	"path"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

func (h *EventHandlers) HandleTranscodeDashJobCompleted(ctx context.Context, ev *events.Event) error {
	var payload messages.JobCompletionPayload
	if err := unmarshalEventData(h.logger, ev, &payload); err != nil {
		return err
	}
	if !payload.Success {
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
	formatVO, err := valueobjects.NewVideoFormat(string(valueobjects.VideoFormatDASH))
	if err != nil {
		return err
	}
	cdnPrefix, playURL := h.cdn.BuildPlayURL(payload.Key)
	si, _ := valueobjects.NewStreamInfo(nil, &cdnPrefix, &playURL)

	cmd := commands.AddVideoCommand{
		AssetID:         *assetIDVO,
		Label:           path.Base(payload.Key),
		Format:          formatVO,
		StorageLocation: *s3Obj,
		ContentType:     payload.ContentType,
		Codec:           payload.VideoCodec,
		VideoCodec:      payload.VideoCodec,
		AudioCodec:      payload.AudioCodec,
		FrameRate:       payload.FrameRate,
		AudioChannels:   payload.AudioChannels,
		AudioSampleRate: payload.AudioSampleRate,
		Duration:        payload.Duration,
		Bitrate:         payload.Bitrate,
		Width:           payload.Width,
		Height:          payload.Height,
		Size:            payload.Size,
		StreamInfo:      si,
	}
	return h.appService.AddVideo(ctx, cmd)
}
