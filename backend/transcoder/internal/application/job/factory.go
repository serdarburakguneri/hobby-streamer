package job

import (
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

type JobFactory struct {
	config config.ServiceConfig
}

func NewJobFactory(config config.ServiceConfig) *JobFactory {
	return &JobFactory{config: config}
}

func (f *JobFactory) CreateJob(payload messages.JobPayload) (*entity.Job, error) {
	assetIDVO, err := valueobjects.NewAssetID(payload.AssetID)
	if err != nil {
		return nil, errors.NewValidationError("invalid asset ID", err)
	}

	videoIDVO, err := valueobjects.NewVideoID(payload.VideoID)
	if err != nil {
		return nil, errors.NewValidationError("invalid video ID", err)
	}

	switch payload.JobType {
	case string(valueobjects.JobTypeAnalyze):
		return f.createAnalyzeJob(*assetIDVO, *videoIDVO, payload)
	case string(valueobjects.JobTypeTranscode):
		return f.createTranscodeJob(*assetIDVO, *videoIDVO, payload)
	default:
		return nil, errors.NewValidationError(fmt.Sprintf("unsupported job type: %s", payload.JobType), nil)
	}
}

func (f *JobFactory) createAnalyzeJob(assetID valueobjects.AssetID, videoID valueobjects.VideoID, payload messages.JobPayload) (*entity.Job, error) {
	return entity.NewAnalyzeJob(assetID, videoID, payload.Input), nil
}

func (f *JobFactory) createTranscodeJob(assetID valueobjects.AssetID, videoID valueobjects.VideoID, payload messages.JobPayload) (*entity.Job, error) {
	format, outputKey, err := f.determineFormatAndOutputKey(payload)
	if err != nil {
		return nil, err
	}

	outputBucket := payload.OutputBucket
	if outputBucket == "" {
		if comp := f.config.GetComponent("s3"); comp != nil {
			if compMap, ok := comp.(map[string]interface{}); ok {
				if bucket, ok2 := compMap["default_output_bucket"].(string); ok2 {
					outputBucket = bucket
				}
			}
		}
	}
	finalOutputKey := payload.OutputKey
	if finalOutputKey == "" {
		finalOutputKey = outputKey
	}

	outputPath := fmt.Sprintf("s3://%s/%s", outputBucket, finalOutputKey)
	return entity.NewTranscodeJob(assetID, videoID, payload.Input, outputPath, payload.Quality, format), nil
}

func (f *JobFactory) determineFormatAndOutputKey(payload messages.JobPayload) (valueobjects.JobFormat, string, error) {
	switch payload.Format {
	case string(valueobjects.JobFormatHLS):
		outputKey := fmt.Sprintf("%s/hls/%s/playlist.m3u8", payload.AssetID, payload.Quality)
		return valueobjects.JobFormatHLS, outputKey, nil
	case string(valueobjects.JobFormatDASH):
		outputKey := fmt.Sprintf("%s/dash/%s/manifest.mpd", payload.AssetID, payload.Quality)
		return valueobjects.JobFormatDASH, outputKey, nil
	default:
		return "", "", errors.NewValidationError(fmt.Sprintf("unsupported format: %s", payload.Format), nil)
	}
}
