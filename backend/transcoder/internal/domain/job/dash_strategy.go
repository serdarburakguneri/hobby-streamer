package job

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
)

type DASHTranscoder struct {
	s3Client *s3.Client
}

func NewDASHTranscoder(s3Client *s3.Client) *DASHTranscoder {
	return &DASHTranscoder{s3Client: s3Client}
}

func (d *DASHTranscoder) Transcode(ctx context.Context, job *Job, localPath, outputDir string) (string, error) {
	outputPath := filepath.Join(outputDir, "playlist.mpd")
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
	retryErr := errors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return "", errors.NewInternalError("DASH transcoding failed", retryErr)
	}

	if job.Type() == JobTypeTranscode && d.s3Client != nil && strings.HasPrefix(job.Output(), "s3://") {
		err := uploadToS3(ctx, d.s3Client, outputDir, job.Output())
		if err != nil {
			return "", errors.NewExternalError("failed to upload DASH output to S3", err)
		}
	}

	return outputPath, nil
}

func (d *DASHTranscoder) ExtractMetadata(ctx context.Context, filePath string, job *Job) (*TranscodeMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.NewInternalError("failed to get file info", err)
	}
	metadata := &TranscodeMetadata{
		Size:        fileInfo.Size(),
		Format:      string(JobFormatDASH),
		ContentType: "application/dash+xml",
	}
	manifest, err := ioutil.ReadFile(filePath)
	if err == nil {
		type Segment struct {
			URL      string `xml:"media,attr"`
			Duration string `xml:"d,attr"`
		}
		type Representation struct {
			ID              string `xml:"id,attr"`
			Codecs          string `xml:"codecs,attr"`
			BaseURL         string `xml:"BaseURL"`
			SegmentTemplate struct {
				Media    string `xml:"media,attr"`
				Duration string `xml:"duration,attr"`
			} `xml:"SegmentTemplate"`
		}
		type MPD struct {
			XMLName xml.Name `xml:"MPD"`
			Period  struct {
				AdaptationSets []struct {
					Representations []Representation `xml:"Representation"`
				} `xml:"AdaptationSet"`
			} `xml:"Period"`
		}
		var mpd MPD
		if xml.Unmarshal(manifest, &mpd) == nil {
			var segments []string
			var segmentCount int
			var totalDuration float64
			for _, as := range mpd.Period.AdaptationSets {
				for _, rep := range as.Representations {
					if rep.Codecs != "" {
						metadata.VideoCodec = rep.Codecs
					}
					if rep.SegmentTemplate.Media != "" && rep.SegmentTemplate.Duration != "" {
						segmentCount = 0
						for i := 1; i <= 10; i++ {
							seg := strings.Replace(rep.SegmentTemplate.Media, "$Number%05d$", fmt.Sprintf("%05d", i), 1)
							segments = append(segments, seg)
							segmentCount++
						}
						dur, _ := strconv.ParseFloat(rep.SegmentTemplate.Duration, 64)
						metadata.AvgSegmentDuration = dur
						totalDuration += float64(segmentCount) * dur
					}
				}
			}
			metadata.SegmentCount = segmentCount
			metadata.Segments = segments
			metadata.Duration = totalDuration
		}
	}

	if job != nil && job.Input() != "" {
		cmd := exec.CommandContext(ctx, "ffprobe",
			"-v", "quiet",
			"-print_format", "json",
			"-show_format",
			"-show_streams",
			job.Input())
		if out, err := cmd.Output(); err == nil {
			var probeResult struct {
				Format struct {
					Duration string `json:"duration"`
					BitRate  string `json:"bit_rate"`
				} `json:"format"`
				Streams []struct {
					CodecType string `json:"codec_type"`
					CodecName string `json:"codec_name"`
					Width     int    `json:"width"`
					Height    int    `json:"height"`
				} `json:"streams"`
			}
			if json.Unmarshal(out, &probeResult) == nil {
				for _, stream := range probeResult.Streams {
					if stream.CodecType == "video" {
						metadata.Width = stream.Width
						metadata.Height = stream.Height
						metadata.Codec = stream.CodecName
						break
					}
				}
				if probeResult.Format.Duration != "" {
					if duration, err := strconv.ParseFloat(probeResult.Format.Duration, 64); err == nil {
						metadata.Duration = duration
					}
				}
				if probeResult.Format.BitRate != "" {
					if bitrate, err := strconv.Atoi(probeResult.Format.BitRate); err == nil {
						metadata.Bitrate = bitrate
					}
				}
			}
		}
	}

	if job != nil && strings.HasPrefix(job.Output(), "s3://") {
		parts := strings.SplitN(job.Output()[5:], "/", 2)
		if len(parts) == 2 {
			metadata.OutputURL = "s3://" + parts[0] + "/" + parts[1]
			metadata.Bucket = parts[0]
			metadata.Key = parts[1]
			metadata.Format = string(JobFormatDASH)
		}
	}
	return metadata, nil
}

func (d *DASHTranscoder) ValidateOutput(job *Job) error {
	if !strings.HasPrefix(job.Output(), "s3://") {
		return errors.NewValidationError("output must be an S3 path", nil)
	}
	parts := strings.SplitN(job.Output()[5:], "/", 2)
	if len(parts) != 2 {
		return errors.NewValidationError("invalid S3 path: "+job.Output(), nil)
	}
	return nil
}

func (d *DASHTranscoder) ValidateInput(ctx context.Context, job *Job) error {
	return nil
}
