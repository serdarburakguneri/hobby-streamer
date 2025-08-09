package transcode

import (
	"context"
	"fmt"
	"path"

	appasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	assetCommands "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	assetQueries "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/queries"
	apppipeline "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/pipeline"
	assetvo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, ev *events.Event) error
}

type Service struct {
	assetCmd  *appasset.CommandService
	assetQry  *appasset.QueryService
	publisher Publisher
	pipeline  *apppipeline.Service
}

func NewService(assetCmd *appasset.CommandService, assetQry *appasset.QueryService, publisher Publisher, pipeline *apppipeline.Service) *Service {
	return &Service{assetCmd: assetCmd, assetQry: assetQry, publisher: publisher, pipeline: pipeline}
}

func (s *Service) RequestTranscode(ctx context.Context, assetID, videoID, format string) error {
	a, err := s.assetQry.GetAsset(ctx, assetQueries.GetAssetQuery{ID: assetID})
	if err != nil || a == nil {
		return fmt.Errorf("asset not found")
	}
	var inputURL, bucket string
	for _, v := range a.Videos() {
		if v.ID().Value() == videoID {
			inputURL = v.StorageLocation().URL()
			bucket = v.StorageLocation().Bucket()
			break
		}
	}
	if inputURL == "" || bucket == "" {
		return fmt.Errorf("video input not found")
	}
	idVO, err := assetvo.NewAssetID(assetID)
	if err != nil {
		return err
	}
	fmtVO, err := assetvo.NewVideoFormat(format)
	if err != nil {
		return err
	}
	var fileName, contentType string
	switch format {
	case assetvo.VideoFormatHLS.Value():
		fileName = "playlist.m3u8"
		contentType = "application/x-mpegURL"
	case assetvo.VideoFormatDASH.Value():
		fileName = "manifest.mpd"
		contentType = "application/dash+xml"
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	outKey := path.Join(assetID, videoID, format, "main", fileName)
	s3Obj, _ := assetvo.NewS3Object(bucket, outKey, "")
	statusTranscoding := assetvo.VideoStatusTranscoding
	if _, _, err := s.assetCmd.UpsertVideo(ctx, assetCommands.UpsertVideoCommand{
		AssetID:         *idVO,
		Label:           fileName,
		Format:          fmtVO,
		StorageLocation: *s3Obj,
		ContentType:     contentType,
		InitialStatus:   &statusTranscoding,
	}); err != nil {
		return err
	}

	corr := events.BuildJobCorrelationID(assetID, videoID, "transcode", format, "main")
	evt := events.NewJobTranscodeRequestedEvent(assetID, videoID, inputURL, format, bucket, outKey)
	evt.SetSource("asset-manager").SetEventVersion("1").SetCorrelationID(corr)
	topic := map[string]string{
		assetvo.VideoFormatHLS.Value():  events.HLSJobRequestedTopic,
		assetvo.VideoFormatDASH.Value(): events.DASHJobRequestedTopic,
	}[format]
	if err := s.publisher.Publish(ctx, topic, evt); err != nil {
		return err
	}
	if s.pipeline != nil {
		_ = s.pipeline.MarkRequested(ctx, assetID, videoID, format, corr, corr)
	}
	return nil
}
