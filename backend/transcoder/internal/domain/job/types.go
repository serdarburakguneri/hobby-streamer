package job

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type VideoMetadata struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Duration    float64 `json:"duration"`
	Bitrate     int     `json:"bitrate"`
	Codec       string  `json:"codec"`
	Size        int64   `json:"size"`
	ContentType string  `json:"contentType"`
}

type TranscodeMetadata struct {
	OutputURL          string   `json:"outputUrl"`
	Bucket             string   `json:"bucket"`
	Key                string   `json:"key"`
	Duration           float64  `json:"duration"`
	Bitrate            int      `json:"bitrate"`
	Size               int64    `json:"size"`
	ContentType        string   `json:"contentType"`
	Format             string   `json:"format"`
	SegmentCount       int      `json:"segmentCount,omitempty"`
	VideoCodec         string   `json:"videoCodec,omitempty"`
	AudioCodec         string   `json:"audioCodec,omitempty"`
	AvgSegmentDuration float64  `json:"avgSegmentDuration,omitempty"`
	Segments           []string `json:"segments,omitempty"`
}

func generateJobID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-only ID if random generation fails
		return time.Now().Format("20060102150405") + "-fallback"
	}
	return time.Now().Format("20060102150405") + "-" + hex.EncodeToString(bytes)
}
