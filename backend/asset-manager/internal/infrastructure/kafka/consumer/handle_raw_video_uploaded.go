package consumer

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
)

func (h *EventHandlers) HandleRawVideoUploaded(ctx context.Context, ev *events.Event) error {
	var payload RawVideoUploadedEvent
	if err := unmarshalEventData(h.logger, ev, &payload); err != nil {
		return err
	}
	assetIDVO, err := valueobjects.NewAssetID(payload.AssetID)
	if err != nil {
		return err
	}
	s3Obj, err := valueobjects.NewS3ObjectFromURL(payload.StorageLocation)
	if err != nil {
		return err
	}
	formatVO, err := valueobjects.NewVideoFormat(string(valueobjects.VideoFormatRaw))
	if err != nil {
		return err
	}
	cmd := commands.AddVideoCommand{
		AssetID:         *assetIDVO,
		Label:           payload.Filename,
		Format:          formatVO,
		StorageLocation: *s3Obj,
		ContentType:     payload.ContentType,
		Codec:           payload.Codec,
		VideoCodec:      "",
		AudioCodec:      "",
		FrameRate:       "",
		AudioChannels:   0,
		AudioSampleRate: 0,
		Duration:        payload.Duration,
		Bitrate:         payload.Bitrate,
		Width:           payload.Width,
		Height:          payload.Height,
		Size:            payload.Size,
		StreamInfo:      nil,
	}
	if err := h.appService.AddVideo(ctx, cmd); err != nil {
		return err
	}
	evt := events.NewJobAnalyzeRequestedEvent(payload.AssetID, payload.VideoID, payload.StorageLocation)
	return h.producer.SendEvent(ctx, events.AnalyzeJobRequestedTopic, evt)
}
