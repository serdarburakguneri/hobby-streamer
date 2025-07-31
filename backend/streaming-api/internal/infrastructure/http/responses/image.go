package responses

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
)

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

func convertImagesToResponse(images []entity.Image) []ImageResponse {
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

func convertImageToResponse(img *entity.Image) *ImageResponse {
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
