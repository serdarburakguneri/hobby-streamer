package asset

type Asset struct {
	ID          int              `json:"id"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt,omitempty"`
	Status      string           `json:"status"`
	Title       *string          `json:"title,omitempty"`
	Description *string          `json:"description,omitempty"`
	Type        *string          `json:"type,omitempty"`
	Category    *string          `json:"category,omitempty"`
	Genres      []string         `json:"genres,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Credits     []Credit         `json:"credits,omitempty"`
	PublishRule *PublishRule     `json:"publishRule,omitempty"`
	Attributes  AttributesMap    `json:"attributes,omitempty"`
	Videos      map[string]Video `json:"videos,omitempty"`
	Images      []Image          `json:"images,omitempty"`
	OwnerID     *string          `json:"ownerId,omitempty"`
}

type Video struct {
	FileName    string           `json:"fileName"`
	ContentType *string          `json:"contentType,omitempty"`
	Duration    *int             `json:"duration,omitempty"`
	Resolution  *string          `json:"resolution,omitempty"`
	Storage     StorageLocations `json:"storage"`
	Stream      *StreamInfo      `json:"stream,omitempty"`
}

type Image struct {
	FileName    string   `json:"fileName"`
	ContentType *string  `json:"contentType,omitempty"`
	Storage     S3Object `json:"storage"`
	AltText     *string  `json:"altText,omitempty"`
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

type Credit struct {
	Name     string  `json:"name"`
	Role     string  `json:"role"`
	Position *string `json:"position,omitempty"`
}

type PublishRule struct {
	PublishDate *string  `json:"publishDate,omitempty"` // RFC3339 timestamp
	ExpireDate  *string  `json:"expireDate,omitempty"`  // RFC3339 timestamp
	RegionLock  []string `json:"regionLock,omitempty"`
}

type AttributesMap map[string]interface{}

// DTOs for API input

type AssetCreateDTO struct {
	Title       *string       `json:"title"`
	Description *string       `json:"description,omitempty"`
	Type        *string       `json:"type"`
	Category    *string       `json:"category,omitempty"`
	Genres      []string      `json:"genres,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Credits     []Credit      `json:"credits,omitempty"`
	PublishRule *PublishRule  `json:"publishRule,omitempty"`
	Attributes  AttributesMap `json:"attributes,omitempty"`
}

type AssetUpdateDTO struct {
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Type        *string       `json:"type,omitempty"`
	Category    *string       `json:"category,omitempty"`
	Genres      []string      `json:"genres,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Credits     []Credit      `json:"credits,omitempty"`
	PublishRule *PublishRule  `json:"publishRule,omitempty"`
	Attributes  AttributesMap `json:"attributes,omitempty"`
	Status      *string       `json:"status,omitempty"`
}
