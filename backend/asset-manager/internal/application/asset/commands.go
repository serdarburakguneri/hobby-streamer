package asset

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
)

type CreateAssetCommand struct {
	Slug        string                 `json:"slug"`
	Title       *string                `json:"title,omitempty"`
	Description *string                `json:"description,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Genre       *string                `json:"genre,omitempty"`
	Genres      []string               `json:"genres,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	OwnerID     *string                `json:"ownerId,omitempty"`
	ParentID    *string                `json:"parentId,omitempty"`
	PublishRule *asset.PublishRule     `json:"publishRule,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func (c CreateAssetCommand) Validate() error {
	if c.Slug == "" {
		return asset.ErrInvalidSlug
	}

	if c.Title != nil && *c.Title == "" {
		return asset.ErrInvalidTitle
	}

	if c.Type != nil && *c.Type == "" {
		return asset.ErrInvalidAssetType
	}

	if c.OwnerID != nil && *c.OwnerID == "" {
		return asset.ErrInvalidOwnerID
	}

	return nil
}

func (c CreateAssetCommand) ToDomainValueObjects() (*asset.Slug, *asset.Title, *asset.Description, *asset.AssetType, *asset.Genre, *asset.Genres, *asset.Tags, *asset.OwnerID, *asset.AssetID, error) {
	slug, err := asset.NewSlug(c.Slug)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	var title *asset.Title
	if c.Title != nil {
		title, err = asset.NewTitle(*c.Title)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var description *asset.Description
	if c.Description != nil {
		description, err = asset.NewDescription(*c.Description)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var assetType *asset.AssetType
	if c.Type != nil {
		assetType, err = asset.NewAssetType(*c.Type)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var genre *asset.Genre
	if c.Genre != nil {
		genre, err = asset.NewGenre(*c.Genre)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var genres *asset.Genres
	if c.Genres != nil {
		genres, err = asset.NewGenres(c.Genres)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var tags *asset.Tags
	if c.Tags != nil {
		tags, err = asset.NewTags(c.Tags)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var ownerID *asset.OwnerID
	if c.OwnerID != nil {
		ownerID, err = asset.NewOwnerID(*c.OwnerID)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var parentID *asset.AssetID
	if c.ParentID != nil {
		parentID, err = asset.NewAssetID(*c.ParentID)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	return slug, title, description, assetType, genre, genres, tags, ownerID, parentID, nil
}

type UpdateAssetCommand struct {
	ID          string                 `json:"id"`
	Title       *string                `json:"title,omitempty"`
	Description *string                `json:"description,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Genre       *string                `json:"genre,omitempty"`
	Genres      []string               `json:"genres,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	OwnerID     *string                `json:"ownerId,omitempty"`
	ParentID    *string                `json:"parentId,omitempty"`
	PublishRule *asset.PublishRule     `json:"publishRule,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func (c UpdateAssetCommand) Validate() error {
	if c.ID == "" {
		return asset.ErrInvalidAssetID
	}

	if c.Title != nil && *c.Title == "" {
		return asset.ErrInvalidTitle
	}

	if c.Type != nil && *c.Type == "" {
		return asset.ErrInvalidAssetType
	}

	if c.OwnerID != nil && *c.OwnerID == "" {
		return asset.ErrInvalidOwnerID
	}

	return nil
}

func (c UpdateAssetCommand) ToDomainValueObjects() (*asset.AssetID, *asset.Title, *asset.Description, *asset.AssetType, *asset.Genre, *asset.Genres, *asset.Tags, *asset.OwnerID, *asset.AssetID, error) {
	assetID, err := asset.NewAssetID(c.ID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	var title *asset.Title
	if c.Title != nil {
		title, err = asset.NewTitle(*c.Title)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var description *asset.Description
	if c.Description != nil {
		description, err = asset.NewDescription(*c.Description)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var assetType *asset.AssetType
	if c.Type != nil {
		assetType, err = asset.NewAssetType(*c.Type)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var genre *asset.Genre
	if c.Genre != nil {
		genre, err = asset.NewGenre(*c.Genre)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var genres *asset.Genres
	if c.Genres != nil {
		genres, err = asset.NewGenres(c.Genres)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var tags *asset.Tags
	if c.Tags != nil {
		tags, err = asset.NewTags(c.Tags)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var ownerID *asset.OwnerID
	if c.OwnerID != nil {
		ownerID, err = asset.NewOwnerID(*c.OwnerID)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	var parentID *asset.AssetID
	if c.ParentID != nil {
		parentID, err = asset.NewAssetID(*c.ParentID)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	return assetID, title, description, assetType, genre, genres, tags, ownerID, parentID, nil
}

type PatchAssetCommand struct {
	ID      string               `json:"id"`
	Patches []JSONPatchOperation `json:"patches"`
}

type JSONPatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value,omitempty"`
}

func (c PatchAssetCommand) Validate() error {
	if c.ID == "" {
		return asset.ErrInvalidAssetID
	}
	if len(c.Patches) == 0 {
		return errors.New("no patches provided")
	}
	return nil
}

func (c PatchAssetCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.ID)
}

type DeleteAssetCommand struct {
	ID string `json:"id"`
}

func (c DeleteAssetCommand) Validate() error {
	if c.ID == "" {
		return asset.ErrInvalidAssetID
	}
	return nil
}

func (c DeleteAssetCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.ID)
}

type GetAssetQuery struct {
	ID   string `json:"id,omitempty"`
	Slug string `json:"slug,omitempty"`
}

func (q GetAssetQuery) Validate() error {
	if q.ID == "" && q.Slug == "" {
		return asset.ErrInvalidAssetID
	}
	return nil
}

func (q GetAssetQuery) ToDomainAssetID() (*asset.AssetID, error) {
	if q.ID != "" {
		return asset.NewAssetID(q.ID)
	}
	return nil, nil
}

func (q GetAssetQuery) ToDomainSlug() (*asset.Slug, error) {
	if q.Slug != "" {
		return asset.NewSlug(q.Slug)
	}
	return nil, nil
}

type ListAssetsQuery struct {
	Limit   int                    `json:"limit"`
	LastKey map[string]interface{} `json:"lastKey,omitempty"`
}

func (q ListAssetsQuery) Validate() error {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	return nil
}

type SearchAssetsQuery struct {
	Query   string                 `json:"query"`
	Limit   int                    `json:"limit"`
	LastKey map[string]interface{} `json:"lastKey,omitempty"`
}

func (q SearchAssetsQuery) Validate() error {
	if q.Query == "" {
		return asset.ErrInvalidTitle
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	return nil
}

type AddVideoCommand struct {
	AssetID         string         `json:"assetId"`
	Label           string         `json:"label"`
	Format          string         `json:"format"`
	StorageLocation asset.S3Object `json:"storageLocation"`
}

func (c AddVideoCommand) Validate() error {
	if c.AssetID == "" {
		return asset.ErrInvalidAssetID
	}
	if c.Label == "" {
		return asset.ErrInvalidTitle
	}
	return nil
}

func (c AddVideoCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.AssetID)
}

type RemoveVideoCommand struct {
	AssetID string `json:"assetId"`
	VideoID string `json:"videoId"`
}

func (c RemoveVideoCommand) Validate() error {
	if c.AssetID == "" {
		return asset.ErrInvalidAssetID
	}
	if c.VideoID == "" {
		return asset.ErrInvalidAssetID
	}
	return nil
}

func (c RemoveVideoCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.AssetID)
}

type UpdateVideoStatusCommand struct {
	AssetID string            `json:"assetId"`
	VideoID string            `json:"videoId"`
	Status  asset.VideoStatus `json:"status"`
}

func (c UpdateVideoStatusCommand) Validate() error {
	if c.AssetID == "" {
		return asset.ErrInvalidAssetID
	}
	if c.VideoID == "" {
		return asset.ErrInvalidAssetID
	}
	return nil
}

func (c UpdateVideoStatusCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.AssetID)
}

type AddImageCommand struct {
	AssetID string      `json:"assetId"`
	Image   asset.Image `json:"image"`
}

func (c AddImageCommand) Validate() error {
	if c.AssetID == "" {
		return asset.ErrInvalidAssetID
	}
	return nil
}

func (c AddImageCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.AssetID)
}

type RemoveImageCommand struct {
	AssetID string `json:"assetId"`
	ImageID string `json:"imageId"`
}

func (c RemoveImageCommand) Validate() error {
	if c.AssetID == "" {
		return asset.ErrInvalidAssetID
	}
	if c.ImageID == "" {
		return asset.ErrInvalidAssetID
	}
	return nil
}

func (c RemoveImageCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.AssetID)
}

type PublishAssetCommand struct {
	AssetID     string             `json:"assetId"`
	PublishRule *asset.PublishRule `json:"publishRule"`
}

func (c PublishAssetCommand) Validate() error {
	if c.AssetID == "" {
		return asset.ErrInvalidAssetID
	}
	if c.PublishRule == nil {
		return asset.ErrInvalidPublishDates
	}
	return nil
}

func (c PublishAssetCommand) ToDomainAssetID() (*asset.AssetID, error) {
	return asset.NewAssetID(c.AssetID)
}
