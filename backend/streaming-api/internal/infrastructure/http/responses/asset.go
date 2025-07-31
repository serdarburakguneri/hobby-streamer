package responses

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
)

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

func NewAssetResponse(a *entity.Asset) AssetResponse {
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
		CreatedAt:   a.CreatedAt(),
		UpdatedAt:   a.UpdatedAt(),
		Metadata:    a.Metadata(),
		OwnerID:     ownerID,
		Videos:      convertVideosToResponse(a.Videos()),
		Images:      convertImagesToResponse(a.Images()),
		PublishRule: convertPublishRuleToResponse(a.PublishRule()),
	}
}
