package valueobjects

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
	Width              int      `json:"width,omitempty"`
	Height             int      `json:"height,omitempty"`
	Duration           float64  `json:"duration"`
	Bitrate            int      `json:"bitrate"`
	Codec              string   `json:"codec,omitempty"`
	Size               int64    `json:"size"`
	ContentType        string   `json:"contentType"`
	Format             string   `json:"format"`
	SegmentCount       int      `json:"segmentCount,omitempty"`
	VideoCodec         string   `json:"videoCodec,omitempty"`
	AudioCodec         string   `json:"audioCodec,omitempty"`
	AvgSegmentDuration float64  `json:"avgSegmentDuration,omitempty"`
	Segments           []string `json:"segments,omitempty"`
	FrameRate          string   `json:"frameRate,omitempty"`
	AudioChannels      int      `json:"audioChannels,omitempty"`
	AudioSampleRate    int      `json:"audioSampleRate,omitempty"`
}
