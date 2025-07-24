package http

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
)

type BucketsResponse struct {
	Buckets []BucketResponse `json:"buckets"`
	Count   int              `json:"count"`
}

type BucketResponse struct {
	ID          string          `json:"id"`
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Type        string          `json:"type"`
	Status      *string         `json:"status,omitempty"`
	AssetIDs    []string        `json:"assetIds,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Assets      []AssetResponse `json:"assets,omitempty"`
}

type AssetsResponse struct {
	Assets []AssetResponse `json:"assets"`
	Count  int             `json:"count"`
}

type AssetResponse struct {
	ID          string               `json:"id"`
	Slug        string               `json:"slug"`
	Title       *string              `json:"title,omitempty"`
	Description *string              `json:"description,omitempty"`
	Type        string               `json:"type"`
	Genre       *string              `json:"genre,omitempty"`
	Genres      []string             `json:"genres,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Status      *string              `json:"status,omitempty"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Metadata    *string              `json:"metadata,omitempty"`
	OwnerID     *string              `json:"ownerId,omitempty"`
	Videos      []VideoResponse      `json:"videos,omitempty"`
	Images      []ImageResponse      `json:"images,omitempty"`
	PublishRule *PublishRuleResponse `json:"publishRule,omitempty"`
}

type VideoResponse struct {
	ID                 string                   `json:"id"`
	Type               string                   `json:"type"`
	Format             string                   `json:"format"`
	StorageLocation    S3ObjectResponse         `json:"storageLocation"`
	Width              *int                     `json:"width,omitempty"`
	Height             *int                     `json:"height,omitempty"`
	Duration           *float64                 `json:"duration,omitempty"`
	Bitrate            *int                     `json:"bitrate,omitempty"`
	Codec              *string                  `json:"codec,omitempty"`
	Size               *int                     `json:"size,omitempty"`
	ContentType        *string                  `json:"contentType,omitempty"`
	StreamInfo         *StreamInfoResponse      `json:"streamInfo,omitempty"`
	Metadata           *string                  `json:"metadata,omitempty"`
	Status             *string                  `json:"status,omitempty"`
	Thumbnail          *ImageResponse           `json:"thumbnail,omitempty"`
	CreatedAt          time.Time                `json:"createdAt"`
	UpdatedAt          time.Time                `json:"updatedAt"`
	Quality            *string                  `json:"quality,omitempty"`
	IsReady            bool                     `json:"isReady"`
	IsProcessing       bool                     `json:"isProcessing"`
	IsFailed           bool                     `json:"isFailed"`
	SegmentCount       *int                     `json:"segmentCount,omitempty"`
	VideoCodec         *string                  `json:"videoCodec,omitempty"`
	AudioCodec         *string                  `json:"audioCodec,omitempty"`
	AvgSegmentDuration *float64                 `json:"avgSegmentDuration,omitempty"`
	Segments           []string                 `json:"segments,omitempty"`
	FrameRate          *string                  `json:"frameRate,omitempty"`
	AudioChannels      *int                     `json:"audioChannels,omitempty"`
	AudioSampleRate    *int                     `json:"audioSampleRate,omitempty"`
	TranscodingInfo    *TranscodingInfoResponse `json:"transcodingInfo,omitempty"`
}

type TranscodingInfoResponse struct {
	JobID       *string    `json:"jobId,omitempty"`
	Progress    *float64   `json:"progress,omitempty"`
	OutputURL   *string    `json:"outputUrl,omitempty"`
	Error       *string    `json:"error,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type ImageResponse struct {
	ID              string              `json:"id"`
	FileName        string              `json:"fileName"`
	URL             string              `json:"url"`
	Type            string              `json:"type"`
	StorageLocation *S3ObjectResponse   `json:"storageLocation,omitempty"`
	Width           *int                `json:"width,omitempty"`
	Height          *int                `json:"height,omitempty"`
	Size            *int                `json:"size,omitempty"`
	ContentType     *string             `json:"contentType,omitempty"`
	StreamInfo      *StreamInfoResponse `json:"streamInfo,omitempty"`
	Metadata        *string             `json:"metadata,omitempty"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

type S3ObjectResponse struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type StreamInfoResponse struct {
	DownloadURL *string `json:"downloadUrl,omitempty"`
	CDNPrefix   *string `json:"cdnPrefix,omitempty"`
	URL         *string `json:"url,omitempty"`
}

type PublishRuleResponse struct {
	PublishAt   *time.Time `json:"publishAt,omitempty"`
	UnpublishAt *time.Time `json:"unpublishAt,omitempty"`
	Regions     []string   `json:"regions,omitempty"`
	AgeRating   *string    `json:"ageRating,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewBucketResponse(b *bucket.Bucket) BucketResponse {
	var description *string
	if b.Description() != nil {
		desc := b.Description().Value()
		description = &desc
	}

	var status *string
	if b.Status() != nil {
		statusVal := b.Status().Value()
		status = &statusVal
	}

	var assetIDs []string
	if b.AssetIDs() != nil {
		assetIDs = b.AssetIDs().Values()
	}

	return BucketResponse{
		ID:          b.ID().Value(),
		Key:         b.Key().Value(),
		Name:        b.Name().Value(),
		Description: description,
		Type:        b.Type().Value(),
		Status:      status,
		AssetIDs:    assetIDs,
		CreatedAt:   b.CreatedAt().Value(),
		UpdatedAt:   b.UpdatedAt().Value(),
	}
}

func NewAssetResponse(a *asset.Asset) AssetResponse {
	var title *string
	if a.Title() != nil {
		titleVal := a.Title().Value()
		title = &titleVal
	}

	var description *string
	if a.Description() != nil {
		desc := a.Description().Value()
		description = &desc
	}

	var genre *string
	if a.Genre() != nil {
		genreVal := a.Genre().Value()
		genre = &genreVal
	}

	var genres []string
	if a.Genres() != nil {
		genreValues := a.Genres().Values()
		genres = make([]string, len(genreValues))
		for i, genre := range genreValues {
			genres[i] = genre.Value()
		}
	}

	var tags []string
	if a.Tags() != nil {
		tags = a.Tags().Values()
	}

	var status *string
	if a.Status() != nil {
		statusVal := a.Status().Value()
		status = &statusVal
	}

	var ownerID *string
	if a.OwnerID() != nil {
		ownerIDVal := a.OwnerID().Value()
		ownerID = &ownerIDVal
	}

	return AssetResponse{
		ID:          a.ID().Value(),
		Slug:        a.Slug().Value(),
		Title:       title,
		Description: description,
		Type:        a.Type().Value(),
		Genre:       genre,
		Genres:      genres,
		Tags:        tags,
		Status:      status,
		CreatedAt:   a.CreatedAt().Value(),
		UpdatedAt:   a.UpdatedAt().Value(),
		Metadata:    a.Metadata(),
		OwnerID:     ownerID,
		Videos:      convertVideosToResponse(a.Videos()),
		Images:      convertImagesToResponse(a.Images()),
		PublishRule: convertPublishRuleToResponse(a.PublishRule()),
	}
}

func convertVideosToResponse(videos []asset.Video) []VideoResponse {
	var response []VideoResponse
	for _, v := range videos {
		var videoType string
		if v.Type() != nil {
			videoType = v.Type().Value()
		}

		var format string
		if v.Format() != nil {
			format = v.Format().Value()
		}

		var quality *string
		if v.Quality() != nil {
			q := v.Quality().Value()
			quality = &q
		}

		var transcodingInfo *TranscodingInfoResponse
		if v.TranscodingInfo() != nil {
			tr := v.TranscodingInfo()
			transcodingInfo = &TranscodingInfoResponse{
				JobID:       tr.JobID,
				Progress:    tr.Progress,
				OutputURL:   tr.OutputURL,
				Error:       tr.Error,
				CompletedAt: tr.CompletedAt,
			}
		}

		response = append(response, VideoResponse{
			ID:                 v.ID().Value(),
			Type:               videoType,
			Format:             format,
			StorageLocation:    convertS3ObjectToResponse(v.StorageLocation()),
			Width:              v.Width(),
			Height:             v.Height(),
			Duration:           v.Duration(),
			Bitrate:            v.Bitrate(),
			Codec:              v.Codec(),
			Size:               v.Size(),
			ContentType:        v.ContentType(),
			StreamInfo:         convertStreamInfoToResponse(v.StreamInfo()),
			Metadata:           v.Metadata(),
			Status:             v.Status(),
			Thumbnail:          convertImageToResponse(v.Thumbnail()),
			CreatedAt:          v.CreatedAt(),
			UpdatedAt:          v.UpdatedAt(),
			Quality:            quality,
			IsReady:            v.IsReadyFlag(),
			IsProcessing:       v.IsProcessing(),
			IsFailed:           v.IsFailed(),
			SegmentCount:       v.SegmentCount(),
			VideoCodec:         v.VideoCodec(),
			AudioCodec:         v.AudioCodec(),
			AvgSegmentDuration: v.AvgSegmentDuration(),
			Segments:           v.Segments(),
			FrameRate:          v.FrameRate(),
			AudioChannels:      v.AudioChannels(),
			AudioSampleRate:    v.AudioSampleRate(),
			TranscodingInfo:    transcodingInfo,
		})
	}
	return response
}

func convertImagesToResponse(images []asset.Image) []ImageResponse {
	var response []ImageResponse
	for _, img := range images {
		var imageType string
		if img.Type() != nil {
			imageType = img.Type().Value()
		}

		response = append(response, ImageResponse{
			ID:              img.ID().Value(),
			FileName:        img.FileName().Value(),
			URL:             img.URL(),
			Type:            imageType,
			StorageLocation: convertS3ObjectToResponsePtr(img.StorageLocation()),
			Width:           img.Width(),
			Height:          img.Height(),
			Size:            img.Size(),
			ContentType:     img.ContentType(),
			StreamInfo:      convertStreamInfoToResponse(img.StreamInfo()),
			Metadata:        img.Metadata(),
			CreatedAt:       img.CreatedAt(),
			UpdatedAt:       img.UpdatedAt(),
		})
	}
	return response
}

func convertS3ObjectToResponse(s3Obj asset.S3ObjectValue) S3ObjectResponse {
	return S3ObjectResponse{
		Bucket: s3Obj.Bucket(),
		Key:    s3Obj.Key(),
		URL:    s3Obj.URL(),
	}
}

func convertS3ObjectToResponsePtr(s3Obj *asset.S3ObjectValue) *S3ObjectResponse {
	if s3Obj == nil {
		return nil
	}
	return &S3ObjectResponse{
		Bucket: s3Obj.Bucket(),
		Key:    s3Obj.Key(),
		URL:    s3Obj.URL(),
	}
}

func convertStreamInfoToResponse(streamInfo *asset.StreamInfoValue) *StreamInfoResponse {
	if streamInfo == nil {
		return nil
	}
	return &StreamInfoResponse{
		DownloadURL: streamInfo.DownloadURL(),
		CDNPrefix:   streamInfo.CDNPrefix(),
		URL:         streamInfo.URL(),
	}
}

func convertImageToResponse(img *asset.Image) *ImageResponse {
	if img == nil {
		return nil
	}

	var imageType string
	if img.Type() != nil {
		imageType = img.Type().Value()
	}

	return &ImageResponse{
		ID:              img.ID().Value(),
		FileName:        img.FileName().Value(),
		URL:             img.URL(),
		Type:            imageType,
		StorageLocation: convertS3ObjectToResponsePtr(img.StorageLocation()),
		Width:           img.Width(),
		Height:          img.Height(),
		Size:            img.Size(),
		ContentType:     img.ContentType(),
		StreamInfo:      convertStreamInfoToResponse(img.StreamInfo()),
		Metadata:        img.Metadata(),
		CreatedAt:       img.CreatedAt(),
		UpdatedAt:       img.UpdatedAt(),
	}
}

func convertPublishRuleToResponse(rule *asset.PublishRuleValue) *PublishRuleResponse {
	if rule == nil {
		return nil
	}

	return &PublishRuleResponse{
		PublishAt:   rule.PublishAt(),
		UnpublishAt: rule.UnpublishAt(),
		Regions:     rule.Regions(),
		AgeRating:   rule.AgeRating(),
	}
}
