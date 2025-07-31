package entity

import (
	"errors"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

type Asset struct {
	id          valueobjects.AssetID
	slug        valueobjects.Slug
	title       *valueobjects.Title
	description *valueobjects.Description
	assetType   valueobjects.AssetType
	genre       *valueobjects.Genre
	genres      *valueobjects.Genres
	tags        *valueobjects.Tags
	status      *valueobjects.Status
	createdAt   time.Time
	updatedAt   time.Time
	metadata    *string
	ownerID     *valueobjects.OwnerID
	videos      []Video
	images      []Image
	publishRule *valueobjects.PublishRuleValue
}

func NewAsset(
	id valueobjects.AssetID,
	slug valueobjects.Slug,
	title *valueobjects.Title,
	description *valueobjects.Description,
	assetType valueobjects.AssetType,
	genre *valueobjects.Genre,
	genres *valueobjects.Genres,
	tags *valueobjects.Tags,
	status *valueobjects.Status,
	createdAt time.Time,
	updatedAt time.Time,
	metadata *string,
	ownerID *valueobjects.OwnerID,
	videos []Video,
	images []Image,
	publishRule *valueobjects.PublishRuleValue,
) *Asset {
	return &Asset{
		id:          id,
		slug:        slug,
		title:       title,
		description: description,
		assetType:   assetType,
		genre:       genre,
		genres:      genres,
		tags:        tags,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		metadata:    metadata,
		ownerID:     ownerID,
		videos:      videos,
		images:      images,
		publishRule: publishRule,
	}
}

func (a *Asset) ID() valueobjects.AssetID {
	return a.id
}

func (a *Asset) Slug() valueobjects.Slug {
	return a.slug
}

func (a *Asset) Title() *valueobjects.Title {
	return a.title
}

func (a *Asset) Description() *valueobjects.Description {
	return a.description
}

func (a *Asset) Type() valueobjects.AssetType {
	return a.assetType
}

func (a *Asset) Genre() *valueobjects.Genre {
	return a.genre
}

func (a *Asset) Genres() *valueobjects.Genres {
	return a.genres
}

func (a *Asset) Tags() *valueobjects.Tags {
	return a.tags
}

func (a *Asset) Status() *valueobjects.Status {
	return a.status
}

func (a *Asset) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Asset) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Asset) Metadata() *string {
	return a.metadata
}

func (a *Asset) OwnerID() *valueobjects.OwnerID {
	return a.ownerID
}

func (a *Asset) Videos() []Video {
	return a.videos
}

func (a *Asset) Images() []Image {
	return a.images
}

func (a *Asset) PublishRule() *valueobjects.PublishRuleValue {
	return a.publishRule
}

func (a *Asset) IsPublished() bool {
	if a.publishRule == nil {
		return false
	}

	now := time.Now().UTC()

	if a.publishRule.PublishAt() != nil && now.Before(*a.publishRule.PublishAt()) {
		return false
	}

	if a.publishRule.UnpublishAt() != nil && now.After(*a.publishRule.UnpublishAt()) {
		return false
	}

	return true
}

func (a *Asset) IsReady() bool {
	return len(a.GetReadyVideos()) > 0
}

func (a *Asset) GetReadyVideos() []Video {
	var readyVideos []Video
	for _, video := range a.videos {
		if video.Status() != nil && video.Status().Value() == constants.VideoStatusReady {
			readyVideos = append(readyVideos, video)
		}
	}
	return readyVideos
}

func (a *Asset) GetMainVideo() *Video {
	for _, video := range a.videos {
		if video.Type() != nil && video.Type().Value() == constants.VideoTypeMain {
			return &video
		}
	}
	return nil
}

func (a *Asset) GetThumbnail() *Image {
	for _, image := range a.images {
		if image.Type() != nil && image.Type().Value() == constants.ImageTypeThumbnail {
			return &image
		}
	}
	return nil
}

func (a *Asset) GetPoster() *Image {
	for _, image := range a.images {
		if image.Type() != nil && image.Type().Value() == constants.ImageTypePoster {
			return &image
		}
	}
	return nil
}

func (a *Asset) HasVideo() bool {
	return len(a.videos) > 0
}

func (a *Asset) HasImage() bool {
	return len(a.images) > 0
}

func (a *Asset) IsAccessibleBy(userID string) bool {
	if a.ownerID == nil {
		return true
	}
	return a.ownerID.Value() == userID
}

func (a *Asset) IsAvailableInRegion(region string) bool {
	if a.publishRule == nil {
		return false
	}

	regions := a.publishRule.Regions()
	if len(regions) == 0 {
		return true
	}

	for _, r := range regions {
		if r == region {
			return true
		}
	}
	return false
}

func (a *Asset) IsAgeAppropriate(userAge int) bool {
	if a.publishRule == nil || a.publishRule.AgeRating() == nil {
		return true
	}

	rating := *a.publishRule.AgeRating()

	ageMap := map[string]int{
		constants.AgeRatingG: 0, constants.AgeRatingPG: 0, constants.AgeRatingPG13: 13, constants.AgeRatingR: 17, constants.AgeRatingNC17: 18,
		constants.AgeRatingTVY: 0, constants.AgeRatingTVY7: 7, constants.AgeRatingTVG: 0, constants.AgeRatingTVPG: 0, constants.AgeRatingTV14: 14, constants.AgeRatingTVMA: 17,
	}

	if minAge, exists := ageMap[rating]; exists {
		return userAge >= minAge
	}

	return true
}

var (
	ErrInvalidAssetType = errors.New("invalid asset type")
	ErrInvalidSlug      = errors.New("invalid slug")
)
