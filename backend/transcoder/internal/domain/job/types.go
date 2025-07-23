package job

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
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

type JobID struct {
	value string
}

func NewJobID(value string) (*JobID, error) {
	if value == "" {
		return nil, errors.New("invalid job ID")
	}
	return &JobID{value: value}, nil
}

func (id JobID) Value() string {
	return id.value
}

func (id JobID) Equals(other JobID) bool {
	return id.value == other.value
}

type AssetID struct {
	value string
}

func NewAssetID(value string) (*AssetID, error) {
	if value == "" {
		return nil, errors.New("invalid asset ID")
	}
	return &AssetID{value: value}, nil
}

func (id AssetID) Value() string {
	return id.value
}

func (id AssetID) Equals(other AssetID) bool {
	return id.value == other.value
}

type VideoID struct {
	value string
}

func NewVideoID(value string) (*VideoID, error) {
	if value == "" {
		return nil, errors.New("invalid video ID")
	}
	return &VideoID{value: value}, nil
}

func (id VideoID) Value() string {
	return id.value
}

func (id VideoID) Equals(other VideoID) bool {
	return id.value == other.value
}

func generateJobID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-only ID if random generation fails
		return time.Now().Format("20060102150405") + "-fallback"
	}
	return time.Now().Format("20060102150405") + "-" + hex.EncodeToString(bytes)
}
