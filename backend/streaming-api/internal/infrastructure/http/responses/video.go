package responses

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
)

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

func convertVideosToResponse(videos []entity.Video) []VideoResponse {
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

		var status *string
		if v.Status() != nil {
			s := v.Status().Value()
			status = &s
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
			Status:             status,
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
