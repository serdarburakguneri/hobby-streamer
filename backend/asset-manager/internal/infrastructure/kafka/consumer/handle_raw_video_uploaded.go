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
	initial := valueobjects.VideoStatusReady
	_, _, err = h.appService.UpsertVideo(ctx, commands.UpsertVideoCommand{
		AssetID:         *assetIDVO,
		Label:           payload.Filename,
		Format:          formatVO,
		StorageLocation: *s3Obj,
		ContentType:     payload.ContentType,
		Codec:           payload.Codec,
		Duration:        payload.Duration,
		Bitrate:         payload.Bitrate,
		Width:           payload.Width,
		Height:          payload.Height,
		Size:            payload.Size,
		InitialStatus:   &initial,
	})
	if err != nil {
		return err
	}
	evt := events.NewJobAnalyzeRequestedEvent(payload.AssetID, payload.VideoID, payload.StorageLocation)
	corr := events.BuildJobCorrelationID(payload.AssetID, payload.VideoID, "analyze", "", "main")
	evt.SetSource("asset-manager").SetEventVersion("1").SetCorrelationID(corr).SetCausationID(ev.ID)
	if h.pipeline != nil {
		_ = h.pipeline.MarkRequested(ctx, payload.AssetID, payload.VideoID, "analyze", corr, corr)
	}
	return h.publisher.Publish(ctx, events.AnalyzeJobRequestedTopic, evt)
}
