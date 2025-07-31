package messages

type JobPayload struct {
	JobType      string `json:"jobType"`
	Input        string `json:"input"`
	AssetID      string `json:"assetId"`
	VideoID      string `json:"videoId"`
	Format       string `json:"format,omitempty"`
	Quality      string `json:"quality,omitempty"`
	OutputBucket string `json:"outputBucket,omitempty"`
	OutputKey    string `json:"outputKey,omitempty"`
}

type JobCompletionPayload struct {
	JobType            string   `json:"jobType"`
	AssetID            string   `json:"assetId"`
	VideoID            string   `json:"videoId"`
	Format             string   `json:"format,omitempty"`
	Success            bool     `json:"success"`
	Error              string   `json:"error,omitempty"`
	Width              int      `json:"width,omitempty"`
	Height             int      `json:"height,omitempty"`
	Duration           float64  `json:"duration,omitempty"`
	Bitrate            int      `json:"bitrate,omitempty"`
	Codec              string   `json:"codec,omitempty"`
	Size               int64    `json:"size,omitempty"`
	ContentType        string   `json:"contentType,omitempty"`
	Bucket             string   `json:"bucket,omitempty"`
	Key                string   `json:"key,omitempty"`
	URL                string   `json:"url,omitempty"`
	SegmentCount       int      `json:"segmentCount,omitempty"`
	VideoCodec         string   `json:"videoCodec,omitempty"`
	AudioCodec         string   `json:"audioCodec,omitempty"`
	AvgSegmentDuration float64  `json:"avgSegmentDuration,omitempty"`
	Segments           []string `json:"segments,omitempty"`
	FrameRate          string   `json:"frameRate,omitempty"`
	AudioChannels      int      `json:"audioChannels,omitempty"`
	AudioSampleRate    int      `json:"audioSampleRate,omitempty"`
}

const (
	MessageTypeJob          = "job"
	MessageTypeJobCompleted = "job-completed"
)
