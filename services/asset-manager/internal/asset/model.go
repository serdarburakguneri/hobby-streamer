package asset

type Asset struct {
	Id           int               `json:"id"`
	UploadDate   string            `json:"uploadDate"`
	Status       string            `json:"status"`
	Title        *string           `json:"title,omitempty"`
	Description  *string           `json:"description,omitempty"`
	ThumbnailURL *string           `json:"thumbnailUrl,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Attributes   AttributesMap     `json:"attributes,omitempty"`
	Variants     map[string]Variant `json:"variants,omitempty"` // e.g., "hd", "trailer", "4k"
}

type Variant struct {
	FileName    string           `json:"fileName"`
	ContentType *string          `json:"contentType,omitempty"`
	Duration    *int             `json:"duration,omitempty"`
	Resolution  *string          `json:"resolution,omitempty"`
	Storage     StorageLocations `json:"storage"`
	Stream      *StreamInfo      `json:"stream,omitempty"`
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

const (
	VariantHD      = "hd"
	Variant4K      = "4k"
	VariantTrailer = "trailer"
)