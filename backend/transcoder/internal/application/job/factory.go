package job

import (
	"bytes"
	"fmt"
	"text/template"

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

func (f *JobFactory) applyTemplate(pattern string, data interface{}) (string, error) {
	tpl, err := template.New("pattern").Parse(pattern)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (f *JobFactory) createTranscodeJob(assetID valueobjects.AssetID, videoID valueobjects.VideoID, payload messages.JobPayload) (*entity.Job, error) {
	comp := f.config.GetComponent("s3").(map[string]interface{})
	bucket := comp["default_output_bucket"].(string)
	sourcePattern := comp["source_prefix_pattern"].(string)
	var filename string
	if payload.OutputKey != "" {
		filename = payload.OutputKey
	}
	sourceKey, err := f.applyTemplate(sourcePattern, map[string]string{
		"AssetID":  assetID.Value(),
		"VideoID":  videoID.Value(),
		"Filename": filename,
	})
	if err != nil {
		return nil, err
	}
	outputPatternKey := ""
	switch payload.Format {
	case string(valueobjects.JobFormatHLS):
		outputPatternKey = comp["hls_output_key_pattern"].(string)
	case string(valueobjects.JobFormatDASH):
		outputPatternKey = comp["dash_output_key_pattern"].(string)
	}
	outputKey, err := f.applyTemplate(outputPatternKey, map[string]string{
		"AssetID": assetID.Value(),
		"VideoID": videoID.Value(),
		"Quality": payload.Quality,
	})
	if err != nil {
		return nil, err
	}
	input := fmt.Sprintf("s3://%s/%s", bucket, sourceKey)
	output := fmt.Sprintf("s3://%s/%s", bucket, outputKey)
	return entity.NewTranscodeJob(assetID, videoID, input, output, payload.Quality, valueobjects.JobFormat(payload.Format)), nil
}

func (f *JobFactory) createAnalyzeJob(assetID valueobjects.AssetID, videoID valueobjects.VideoID, payload messages.JobPayload) (*entity.Job, error) {
	return entity.NewAnalyzeJob(assetID, videoID, payload.Input), nil
}
