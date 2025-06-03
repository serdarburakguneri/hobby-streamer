package asset

type Asset struct {
	Id           int           `json:"id"`
	FileName     string        `json:"fileName"`
	UploadDate   string        `json:"uploadDate"`
	Status       string        `json:"status"`
	ContentType  *string       `json:"contentType,omitempty"`
	Duration     *int          `json:"duration,omitempty"`
	Resolution   *string       `json:"resolution,omitempty"`
	Title        *string       `json:"title,omitempty"`
	Description  *string       `json:"description,omitempty"`
	Storage      *StorageLocations `json:"storage,omitempty"`
	ThumbnailURL *string       `json:"thumbnailUrl,omitempty"`
	Stream       *StreamInfo   `json:"stream,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Attributes   AttributesMap `json:"attributes,omitempty"`
}

type StorageLocations struct {
	Raw        *S3Object `json:"raw,omitempty"`
	Transcoded *S3Object `json:"transcoded,omitempty"`
	Thumbnail  *S3Object `json:"thumbnail,omitempty"`
}

type S3Object struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type StreamInfo struct {
	HLS         *string `json:"hls,omitempty"`
	DASH        *string `json:"dash,omitempty"`
	DownloadURL *string `json:"downloadUrl,omitempty"`
	CdnPrefix   *string `json:"cdnPrefix,omitempty"`
}

type AttributesMap map[string]interface{}