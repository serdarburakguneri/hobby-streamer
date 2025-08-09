package entity

import (
	"errors"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/operations"
)

type Asset struct {
	id          valueobjects.AssetID
	version     int
	slug        valueobjects.Slug
	title       *valueobjects.Title
	description *valueobjects.Description
	assetType   *valueobjects.AssetType
	genre       *valueobjects.Genre
	genres      *valueobjects.Genres
	tags        *valueobjects.Tags
	createdAt   valueobjects.CreatedAt
	updatedAt   valueobjects.UpdatedAt
	ownerID     *valueobjects.OwnerID
	parentID    *valueobjects.AssetID
	parent      *Asset
	children    []Asset
	images      []valueobjects.Image
	videos      map[string]*Video
	credits     []valueobjects.Credit
	publishRule *valueobjects.PublishRule
	metadata    map[string]interface{}
}

func NewAsset(slug valueobjects.Slug, title *valueobjects.Title, assetType *valueobjects.AssetType) (*Asset, error) {
	now := time.Now().UTC()
	assetID, err := valueobjects.NewAssetID(operations.GenerateID())
	if err != nil {
		return nil, err
	}
	createdAt := valueobjects.NewCreatedAt(now)
	updatedAt := valueobjects.NewUpdatedAt(now)
	genres, err := valueobjects.NewGenres([]string{})
	if err != nil {
		return nil, err
	}
	tags, err := valueobjects.NewTags([]string{})
	if err != nil {
		return nil, err
	}

	return &Asset{
		id:        *assetID,
		version:   0,
		slug:      slug,
		title:     title,
		assetType: assetType,
		createdAt: *createdAt,
		updatedAt: *updatedAt,
		videos:    make(map[string]*Video),
		images:    make([]valueobjects.Image, 0),
		credits:   make([]valueobjects.Credit, 0),
		genres:    genres,
		tags:      tags,
		metadata:  make(map[string]interface{}),
	}, nil
}

func ReconstructAsset(
	id valueobjects.AssetID,
	slug valueobjects.Slug,
	title *valueobjects.Title,
	description *valueobjects.Description,
	assetType *valueobjects.AssetType,
	genre *valueobjects.Genre,
	genres *valueobjects.Genres,
	tags *valueobjects.Tags,
	createdAt valueobjects.CreatedAt,
	updatedAt valueobjects.UpdatedAt,
	ownerID *valueobjects.OwnerID,
	parentID *valueobjects.AssetID,
	images []valueobjects.Image,
	videos map[string]*Video,
	credits []valueobjects.Credit,
	publishRule *valueobjects.PublishRule,
	metadata map[string]interface{},
) *Asset {
	return &Asset{
		id:          id,
		version:     0,
		slug:        slug,
		title:       title,
		description: description,
		assetType:   assetType,
		genre:       genre,
		genres:      genres,
		tags:        tags,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		ownerID:     ownerID,
		parentID:    parentID,
		images:      images,
		videos:      videos,
		credits:     credits,
		publishRule: publishRule,
		metadata:    metadata,
	}
}

func (a *Asset) ID() valueobjects.AssetID {
	return a.id
}

func (a *Asset) Version() int     { return a.version }
func (a *Asset) SetVersion(v int) { a.version = v }

func (a *Asset) Slug() valueobjects.Slug {
	return a.slug
}

func (a *Asset) Title() *valueobjects.Title {
	return a.title
}

func (a *Asset) Description() *valueobjects.Description {
	return a.description
}

func (a *Asset) Type() *valueobjects.AssetType {
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

func (a *Asset) CreatedAt() valueobjects.CreatedAt {
	return a.createdAt
}

func (a *Asset) UpdatedAt() valueobjects.UpdatedAt {
	return a.updatedAt
}

func (a *Asset) OwnerID() *valueobjects.OwnerID {
	return a.ownerID
}

func (a *Asset) ParentID() *valueobjects.AssetID {
	return a.parentID
}

func (a *Asset) Parent() *Asset {
	return a.parent
}

func (a *Asset) Children() []Asset {
	return a.children
}

func (a *Asset) Images() []valueobjects.Image {
	return a.images
}

func (a *Asset) Videos() map[string]*Video {
	return a.videos
}

func (a *Asset) Credits() []valueobjects.Credit {
	return a.credits
}

func (a *Asset) PublishRule() *valueobjects.PublishRule {
	return a.publishRule
}

func (a *Asset) Metadata() map[string]interface{} {
	return a.metadata
}

func (a *Asset) UpdateTitle(newTitle *valueobjects.Title) {
	a.title = newTitle
	a.touch()
}

func (a *Asset) UpdateDescription(newDescription *valueobjects.Description) {
	a.description = newDescription
	a.touch()
}

func (a *Asset) UpdateGenre(newGenre *valueobjects.Genre) {
	a.genre = newGenre
	a.touch()
}

func (a *Asset) UpdateGenres(newGenres *valueobjects.Genres) {
	a.genres = newGenres
	a.touch()
}

func (a *Asset) UpdateTags(newTags *valueobjects.Tags) {
	a.tags = newTags
	a.touch()
}

func (a *Asset) Status() string {
	return "draft"
}

func (a *Asset) CanUpdateTitle() bool {
	return true
}

func (a *Asset) CanUpdateDescription() bool {
	return true
}

func (a *Asset) IsReadyForPublishing() bool {
	hasReadyVideos := false
	for _, video := range a.videos {
		if video.Status() == valueobjects.VideoStatus(constants.VideoStatusReady) {
			hasReadyVideos = true
			break
		}
	}

	return hasReadyVideos
}

func (a *Asset) findVideoByLabelAndFormat(label string, format valueobjects.VideoFormat) (string, *Video) {
	for id, v := range a.videos {
		if v.Label().Value() == label && v.Format().Equals(format) {
			return id, v
		}
	}
	return "", nil
}

func (a *Asset) hasMainForFormat(format valueobjects.VideoFormat) bool {
	for _, v := range a.videos {
		if v.Type().IsMain() && v.Format().Equals(format) {
			return true
		}
	}
	return false
}

func (a *Asset) UpsertVideo(
	label string,
	format *valueobjects.VideoFormat,
	storageLocation valueobjects.S3Object,
	width, height int,
	duration float64,
	bitrate int,
	codec string,
	size int64,
	contentType string,
	videoCodec, audioCodec string,
	frameRate string,
	audioChannels, audioSampleRate int,
	streamInfo *valueobjects.StreamInfo,
	initialStatus *valueobjects.VideoStatus,
) (*Video, error) {
	if format == nil {
		return nil, errors.New("format cannot be nil")
	}
	if _, existing := a.findVideoByLabelAndFormat(label, *format); existing != nil {
		existing.UpdateStorageLocation(storageLocation)
		if streamInfo != nil {
			existing.SetStreamInfo(streamInfo)
		}
		ct, err := valueobjects.NewContentType(contentType)
		if err != nil {
			return nil, err
		}
		existing.UpdateMediaInfo(*valueobjects.NewMediaInfo(
			width, height, duration, bitrate, codec, size, *ct, videoCodec, audioCodec, frameRate, audioChannels, audioSampleRate,
		))
		if initialStatus != nil {
			existing.UpdateStatus(*initialStatus)
		}
		a.touch()
		return existing, nil
	}

	if a.hasMainForFormat(*format) {
		return nil, errors.New("a MAIN video for this format already exists")
	}

	video, err := NewVideo(
		label,
		format,
		storageLocation,
		width,
		height,
		duration,
		bitrate,
		codec,
		size,
		contentType,
		videoCodec,
		audioCodec,
		frameRate,
		audioChannels,
		audioSampleRate,
		streamInfo,
	)
	if err != nil {
		return nil, err
	}
	a.videos[video.ID().Value()] = video
	if initialStatus != nil {
		video.UpdateStatus(*initialStatus)
	}
	a.touch()
	return video, nil
}

func (a *Asset) RemoveVideo(videoID string) error {
	if _, exists := a.videos[videoID]; exists {
		delete(a.videos, videoID)
		a.touch()
		return nil
	}
	return errors.New("video not found")
}

func (a *Asset) UpdateVideoStatus(videoID string, status valueobjects.VideoStatus) error {
	if video, exists := a.videos[videoID]; exists {
		video.UpdateStatus(status)
		a.touch()
		return nil
	}
	return errors.New("video not found")
}

func (a *Asset) UpdateVideoMediaInfo(videoID string, transcodingInfo valueobjects.TranscodingInfo) error {
	if video, exists := a.videos[videoID]; exists {
		video.UpdateMediaInfo(transcodingInfo)
		video.UpdateStatus(valueobjects.VideoStatus("ready"))
		a.touch()
		return nil
	}
	return errors.New("video not found")
}

func (a *Asset) AddImage(image valueobjects.Image) {
	a.images = append(a.images, image)
	a.touch()
}

func (a *Asset) RemoveImage(imageID string) error {
	for i, image := range a.images {
		if image.ID().Value() == imageID {
			a.images = append(a.images[:i], a.images[i+1:]...)
			a.touch()
			return nil
		}
	}
	return errors.New("image not found")
}

func (a *Asset) AddCredit(credit valueobjects.Credit) {
	a.credits = append(a.credits, credit)
	a.touch()
}

func (a *Asset) RemoveCredit(personID string) error {
	for i, credit := range a.credits {
		if credit.Name() == personID {
			a.credits = append(a.credits[:i], a.credits[i+1:]...)
			a.touch()
			return nil
		}
	}
	return errors.New("credit not found")
}

func (a *Asset) CanBePublished() bool {
	return a.IsReadyForPublishing()
}

func (a *Asset) SetPublishRule(rule *valueobjects.PublishRule) error {
	a.publishRule = rule
	a.touch()
	return nil
}

func (a *Asset) SetParentID(parentID *valueobjects.AssetID) {
	a.parentID = parentID
	a.touch()
}

func (a *Asset) ValidateForPublishing() error {
	if a.publishRule == nil {
		return errors.New("asset is not ready for publishing")
	}
	return nil
}

func (a *Asset) CalculateMetrics() interface{} {
	return map[string]interface{}{
		"videoCount":  len(a.videos),
		"imageCount":  len(a.images),
		"creditCount": len(a.credits),
	}
}

func (a *Asset) CalculateStorageUsage() interface{} {
	totalSize := int64(0)
	for _, video := range a.videos {
		totalSize += video.Size()
	}
	return map[string]interface{}{
		"totalSize":  totalSize,
		"videoCount": len(a.videos),
	}
}

func (a *Asset) SetOwnerID(ownerID *valueobjects.OwnerID) {
	a.ownerID = ownerID
	a.touch()
}

func (a *Asset) touch() {
	a.updatedAt = *valueobjects.NewUpdatedAt(time.Now().UTC())
}
