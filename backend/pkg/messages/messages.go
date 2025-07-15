package messages

type AnalyzePayload struct {
	Input   string `json:"input"`
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
}

type TranscodePayload struct {
	Input          string `json:"input"`
	AssetID        string `json:"assetId"`
	VideoID        string `json:"videoId"`
	Format         string `json:"format"`
	OutputBucket   string `json:"outputBucket"`
	OutputKey      string `json:"outputKey"`
	OutputFileName string `json:"outputFileName"`
}

type AnalyzeCompletionPayload struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type TranscodeCompletionPayload struct {
	AssetID  string `json:"assetId"`
	VideoID  string `json:"videoId"`
	Format   string `json:"format"`
	Success  bool   `json:"success"`
	Bucket   string `json:"bucket,omitempty"`
	Key      string `json:"key,omitempty"`
	FileName string `json:"fileName,omitempty"`
	URL      string `json:"url,omitempty"`
	Error    string `json:"error,omitempty"`
}

const (
	MessageTypeAnalyze                = "analyze"
	MessageTypeTranscodeHLS           = "transcode-hls"
	MessageTypeTranscodeDASH          = "transcode-dash"
	MessageTypeAnalyzeCompleted       = "analyze-completed"
	MessageTypeTranscodeHLSCompleted  = "transcode-hls-completed"
	MessageTypeTranscodeDASHCompleted = "transcode-dash-completed"
)
