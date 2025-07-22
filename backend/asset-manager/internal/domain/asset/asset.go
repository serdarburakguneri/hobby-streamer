package asset

import (
	"errors"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type Asset struct {
	id          AssetID
	slug        Slug
	title       *Title
	description *Description
	assetType   *AssetType
	genre       *Genre
	genres      *Genres
	tags        *Tags
	createdAt   CreatedAt
	updatedAt   UpdatedAt
	ownerID     *OwnerID
	parentID    *AssetID
	parent      *Asset
	children    []Asset
	images      []Image
	videos      map[string]*Video
	credits     []Credit
	publishRule *PublishRule
	metadata    map[string]interface{}
}

func NewAsset(slug Slug, title *Title, assetType *AssetType) (*Asset, error) {
	now := time.Now().UTC()
	return &Asset{
		id:        AssetID{value: generateID()},
		slug:      slug,
		title:     title,
		assetType: assetType,
		createdAt: CreatedAt{value: now},
		updatedAt: UpdatedAt{value: now},
		videos:    make(map[string]*Video),
		images:    make([]Image, 0),
		credits:   make([]Credit, 0),
		genres:    &Genres{values: make([]Genre, 0)},
		tags:      &Tags{values: make([]Tag, 0)},
		metadata:  make(map[string]interface{}),
	}, nil
}

func ReconstructAsset(
	id AssetID,
	slug Slug,
	title *Title,
	description *Description,
	assetType *AssetType,
	genre *Genre,
	genres *Genres,
	tags *Tags,
	createdAt CreatedAt,
	updatedAt UpdatedAt,
	ownerID *OwnerID,
	parentID *AssetID,
	images []Image,
	videos map[string]*Video,
	credits []Credit,
	publishRule *PublishRule,
	metadata map[string]interface{},
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

func (a *Asset) ID() AssetID {
	return a.id
}

func (a *Asset) Slug() Slug {
	return a.slug
}

func (a *Asset) Title() *Title {
	return a.title
}

func (a *Asset) Description() *Description {
	return a.description
}

func (a *Asset) Type() *AssetType {
	return a.assetType
}

func (a *Asset) Genre() *Genre {
	return a.genre
}

func (a *Asset) Genres() *Genres {
	return a.genres
}

func (a *Asset) Tags() *Tags {
	return a.tags
}

func (a *Asset) CreatedAt() CreatedAt {
	return a.createdAt
}

func (a *Asset) UpdatedAt() UpdatedAt {
	return a.updatedAt
}

func (a *Asset) OwnerID() *OwnerID {
	return a.ownerID
}

func (a *Asset) ParentID() *AssetID {
	return a.parentID
}

func (a *Asset) Parent() *Asset {
	return a.parent
}

func (a *Asset) Children() []Asset {
	return a.children
}

func (a *Asset) Images() []Image {
	return a.images
}

func (a *Asset) Videos() []*Video {
	videos := make([]*Video, 0, len(a.videos))
	for _, video := range a.videos {
		videos = append(videos, video)
	}
	return videos
}

func (a *Asset) Credits() []Credit {
	return a.credits
}

func (a *Asset) PublishRule() *PublishRule {
	return a.publishRule
}

func (a *Asset) Metadata() map[string]interface{} {
	return a.metadata
}

func (a *Asset) Status() string {
	if a.publishRule == nil {
		return constants.AssetStatusDraft
	}

	now := time.Now().UTC()

	if a.publishRule.publishAt == nil {
		return constants.AssetStatusDraft
	}

	if now.Before(*a.publishRule.publishAt) {
		return constants.AssetStatusScheduled
	}

	if a.publishRule.unpublishAt != nil && now.After(*a.publishRule.unpublishAt) {
		return constants.AssetStatusExpired
	}

	return constants.AssetStatusPublished
}

func (a *Asset) CanUpdateTitle() bool {
	return true
}

func (a *Asset) UpdateTitle(title *Title) error {
	if !a.CanUpdateTitle() {
		return ErrCannotUpdateAsset
	}

	a.title = title
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanUpdateDescription() bool {
	return true
}

func (a *Asset) UpdateDescription(description *Description) error {
	a.description = description
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanUpdateType() bool {
	return true
}

func (a *Asset) UpdateType(assetType *AssetType) error {
	if !a.CanUpdateType() {
		return ErrCannotUpdateAsset
	}

	a.assetType = assetType
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) UpdateGenre(genre *Genre) error {
	a.genre = genre
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) SetGenres(genres *Genres) error {
	a.genres = genres
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) SetTags(tags *Tags) error {
	a.tags = tags
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanUpdateOwner() bool {
	return true
}

func (a *Asset) SetOwnerID(ownerID *OwnerID) error {
	if !a.CanUpdateOwner() {
		return ErrCannotUpdateAsset
	}

	a.ownerID = ownerID
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanUpdateParent() bool {
	return true
}

func (a *Asset) SetParentID(parentID *AssetID) error {
	if !a.CanUpdateParent() {
		return ErrCannotUpdateAsset
	}

	if parentID != nil && a.id.Equals(*parentID) {
		return ErrAssetCannotBeOwnParent
	}

	a.parentID = parentID
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanUpdatePublishRule() bool {
	return true
}

func (a *Asset) SetPublishRule(publishRule *PublishRule) error {
	a.publishRule = publishRule
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanUpdateMetadata() bool {
	return true
}

func (a *Asset) SetMetadata(metadata map[string]interface{}) error {
	a.metadata = metadata
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanAddVideo() bool {
	return true
}

func (a *Asset) AddVideo(label string, format *VideoFormat, storageLocation S3Object) (*Video, error) {
	if !a.CanAddVideo() {
		return nil, ErrCannotUpdateAsset
	}

	video := NewVideo(label, format, storageLocation, 0, "", "", 0, nil)
	a.videos[video.ID()] = video
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return video, nil
}

func (a *Asset) CanRemoveVideo() bool {
	return true
}

func (a *Asset) RemoveVideo(videoID string) error {
	if !a.CanRemoveVideo() {
		return ErrCannotRemoveVideo
	}

	if _, exists := a.videos[videoID]; !exists {
		return ErrVideoNotFound
	}

	delete(a.videos, videoID)
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) GetVideo(videoID string) (*Video, error) {
	video, exists := a.videos[videoID]
	if !exists {
		return nil, ErrVideoNotFound
	}
	return video, nil
}

func (a *Asset) CanUpdateVideoStatus() bool {
	return true
}

func (a *Asset) UpdateVideoStatus(videoID string, status VideoStatus) error {
	if !a.CanUpdateVideoStatus() {
		return ErrCannotUpdateAsset
	}

	video, err := a.GetVideo(videoID)
	if err != nil {
		return err
	}

	video.UpdateStatus(status)
	return nil
}

func (a *Asset) CanUpdateVideoTranscodingInfo() bool {
	return true
}

func (a *Asset) UpdateVideoTranscodingInfo(videoID string, transcodingInfo TranscodingInfo) error {
	if !a.CanUpdateVideoTranscodingInfo() {
		return ErrCannotUpdateAsset
	}

	video, err := a.GetVideo(videoID)
	if err != nil {
		return err
	}

	video.UpdateTranscodingInfo(transcodingInfo)
	return nil
}

func (a *Asset) CanAddImage() bool {
	return true
}

func (a *Asset) AddImage(image Image) error {
	if !a.CanAddImage() {
		return ErrCannotUpdateAsset
	}

	for _, existingImage := range a.images {
		if existingImage.ID() == image.ID() {
			return ErrImageAlreadyExists
		}
	}

	a.images = append(a.images, image)
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanRemoveImage() bool {
	return true
}

func (a *Asset) RemoveImage(imageID string) error {
	if !a.CanRemoveImage() {
		return ErrCannotUpdateAsset
	}

	for i, image := range a.images {
		if image.ID() == imageID {
			a.images = append(a.images[:i], a.images[i+1:]...)
			a.updatedAt = UpdatedAt{value: time.Now().UTC()}
			return nil
		}
	}

	return ErrImageNotFound
}

func (a *Asset) CanAddCredit() bool {
	return true
}

func (a *Asset) AddCredit(credit Credit) error {
	if !a.CanAddCredit() {
		return ErrCannotUpdateAsset
	}

	a.credits = append(a.credits, credit)
	a.updatedAt = UpdatedAt{value: time.Now().UTC()}
	return nil
}

func (a *Asset) CanRemoveCredit() bool {
	return true
}

func (a *Asset) RemoveCredit(personID string) error {
	if !a.CanRemoveCredit() {
		return ErrCannotUpdateAsset
	}

	for i, credit := range a.credits {
		if credit.personID != nil && *credit.personID == personID {
			a.credits = append(a.credits[:i], a.credits[i+1:]...)
			a.updatedAt = UpdatedAt{value: time.Now().UTC()}
			return nil
		}
	}

	return ErrCreditNotFound
}

func (a *Asset) IsReadyForPublishing() bool {
	if a.title == nil {
		return false
	}

	if a.assetType == nil {
		return false
	}

	if len(a.videos) == 0 {
		return false
	}

	hasReadyVideo := false
	for _, video := range a.videos {
		if video.IsReady() {
			hasReadyVideo = true
			break
		}
	}

	return hasReadyVideo
}

func (a *Asset) CanBePublished() bool {
	return true
}

func (a *Asset) CanBeUnpublished() bool {
	return true
}

var (
	ErrCannotUpdateAsset      = errors.New("cannot update asset")
	ErrAssetCannotBeOwnParent = errors.New("asset cannot be its own parent")
	ErrVideoNotFound          = errors.New("video not found")
	ErrCannotRemoveVideo      = errors.New("cannot remove video")
	ErrImageNotFound          = errors.New("image not found")
	ErrImageAlreadyExists     = errors.New("image already exists")
	ErrCreditNotFound         = errors.New("credit not found")
)
