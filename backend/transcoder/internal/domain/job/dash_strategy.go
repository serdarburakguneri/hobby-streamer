package job

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type DASHTranscoder struct{}

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
	return outputPath, nil
}

func (d *DASHTranscoder) ExtractMetadata(ctx context.Context, filePath string) (*TranscodeMetadata, error) {
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
	return metadata, nil
}
