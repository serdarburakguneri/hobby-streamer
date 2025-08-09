package transcoding

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	
	baseDir := filepath.Dir(filePath)
	files, _ := ioutil.ReadDir(baseDir)
	var firstSeg string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".m4s") {
			firstSeg = filepath.Join(baseDir, f.Name())
			break
		}
	}
	if firstSeg != "" {
		probe := exec.CommandContext(ctx, "ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", firstSeg)
		if out, err := probe.CombinedOutput(); err == nil {
			type Stream struct {
				CodecType    string `json:"codec_type"`
				CodecName    string `json:"codec_name"`
				Width        int    `json:"width"`
				Height       int    `json:"height"`
				SampleRate   string `json:"sample_rate"`
				Channels     int    `json:"channels"`
				AvgFrameRate string `json:"avg_frame_rate"`
			}
			var pr struct {
				Format struct {
					BitRate string `json:"bit_rate"`
				} `json:"format"`
				Streams []Stream `json:"streams"`
			}
			if json.Unmarshal(out, &pr) == nil {
				for _, s := range pr.Streams {
					if s.CodecType == "video" {
						if s.Width > 0 {
							metadata.Width = s.Width
						}
						if s.Height > 0 {
							metadata.Height = s.Height
						}
						if s.CodecName != "" {
							metadata.VideoCodec = s.CodecName
						}
						if s.AvgFrameRate != "" {
							metadata.FrameRate = s.AvgFrameRate
						}
					} else if s.CodecType == "audio" {
						if s.Channels > 0 {
							metadata.AudioChannels = s.Channels
						}
						if s.SampleRate != "" {
							if sr, err := strconv.Atoi(s.SampleRate); err == nil {
								metadata.AudioSampleRate = sr
							}
						}
						if s.CodecName != "" {
							metadata.AudioCodec = s.CodecName
						}
					}
				}
				if pr.Format.BitRate != "" {
					if b, err := strconv.Atoi(pr.Format.BitRate); err == nil {
						metadata.Bitrate = b
					}
				}
			}
		}
	}
	return metadata, nil
}
