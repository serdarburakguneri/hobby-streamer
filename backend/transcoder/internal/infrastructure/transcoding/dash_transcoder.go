package transcoding

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	resilience "github.com/serdarburakguneri/hobby-streamer/backend/pkg/resilience"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

type DASHTranscoder struct {
	storage job.Storage
}

func NewDASHTranscoder(storage job.Storage) *DASHTranscoder {
	return &DASHTranscoder{storage: storage}
}

func (d *DASHTranscoder) ValidateInput(ctx context.Context, job *entity.Job) error {
	return nil
}

func (d *DASHTranscoder) Transcode(ctx context.Context, job *entity.Job, localPath, outputDir string) (string, error) {
	outputPath := filepath.Join(outputDir, "manifest.mpd")
	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", localPath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-f", "dash",
			"-seg_duration", "10",
			"-use_template", "1",
			"-use_timeline", "1",
			"-init_seg_name", "init-$RepresentationID$.m4s",
			"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s",
			outputPath)
		return cmd.Run()
	}
	if err := resilience.RetryWithBackoff(ctx, retryFunc, 2); err != nil {
		return "", pkgerrors.NewInternalError("DASH transcoding failed", err)
	}
	if job.Type().IsTranscode() && strings.HasPrefix(job.Output(), "s3://") {
		if err := d.storage.Upload(ctx, outputDir, job.Output()); err != nil {
			return "", pkgerrors.NewExternalError("failed to upload DASH output to S3", err)
		}
	}
	return outputPath, nil
}

func (d *DASHTranscoder) ValidateOutput(job *entity.Job) error {
	return nil
}

func (d *DASHTranscoder) ExtractMetadata(ctx context.Context, filePath string, job *entity.Job) (*valueobjects.TranscodeMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to stat DASH manifest", err)
	}
	outputURL := job.Output()
	var bucket, key string
	if strings.HasPrefix(outputURL, "s3://") {
		parts := strings.SplitN(outputURL[5:], "/", 2)
		if len(parts) == 2 {
			bucket = parts[0]
			key = parts[1]
		}
	}
	metadata := &valueobjects.TranscodeMetadata{
		OutputURL:   outputURL,
		Bucket:      bucket,
		Key:         key,
		Size:        fileInfo.Size(),
		Format:      valueobjects.JobFormatDASH.String(),
		ContentType: "application/dash+xml",
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return metadata, nil
	}
	type Representation struct {
		Codecs string `xml:"codecs,attr"`
	}
	type MPD struct {
		Period struct {
			AdaptationSet struct {
				Representations []Representation `xml:"Representation"`
			} `xml:"AdaptationSet"`
		} `xml:"Period"`
	}
	var mpd MPD
	if err := xml.Unmarshal(data, &mpd); err == nil && len(mpd.Period.AdaptationSet.Representations) > 0 {
		codecs := mpd.Period.AdaptationSet.Representations[0].Codecs
		metadata.VideoCodec = codecs
		metadata.Codec = codecs
	}
	return metadata, nil
}
