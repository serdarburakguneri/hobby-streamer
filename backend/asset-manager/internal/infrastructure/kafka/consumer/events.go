package consumer

type TranscoderJobCompletedEvent struct {
	JobID        string                 `json:"jobId"`
	JobType      string                 `json:"jobType"`
	AssetID      string                 `json:"assetId"`
	VideoID      string                 `json:"videoId"`
	Success      bool                   `json:"success"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
	CompletedAt  string                 `json:"completedAt"`
}

type RawVideoUploadedEvent struct {
	AssetID         string  `json:"assetId"`
	VideoID         string  `json:"videoId"`
	StorageLocation string  `json:"storageLocation"`
	Filename        string  `json:"filename"`
	Size            int64   `json:"size"`
	ContentType     string  `json:"contentType"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Duration        float64 `json:"duration"`
	Bitrate         int     `json:"bitrate"`
	Codec           string  `json:"codec"`
}

type AnalyzeJobCompletedEvent struct {
	JobID        string                 `json:"jobId"`
	JobType      string                 `json:"jobType"`
	AssetID      string                 `json:"assetId"`
	VideoID      string                 `json:"videoId"`
	Success      bool                   `json:"success"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
	CompletedAt  string                 `json:"completedAt"`
}
