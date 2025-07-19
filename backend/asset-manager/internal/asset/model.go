package asset

import (
	"time"
)

type Asset struct {
	ID          string                 `json:"id"`
	Slug        string                 `json:"slug"`
	Title       *string                `json:"title,omitempty"`
	Description *string                `json:"description,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Genre       *string                `json:"genre,omitempty"`
	Genres      []string               `json:"genres,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	OwnerID     *string                `json:"ownerId,omitempty"`
	ParentID    *string                `json:"parentId,omitempty"`
	Parent      *Asset                 `json:"parent,omitempty"`
	Children    []Asset                `json:"children,omitempty"`
	Images      []Image                `json:"images,omitempty"`
	Videos      []Video                `json:"videos,omitempty"`
	Credits     []Credit               `json:"credits,omitempty"`
	PublishRule *PublishRule           `json:"publishRule,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type VideoType string

type VideoFormat string

type Video struct {
	ID              string            `json:"id"`
	Type            VideoType         `json:"type"`
	Format          VideoFormat       `json:"format"`
	StorageLocation S3Object          `json:"storageLocation"`
	Width           int               `json:"width,omitempty"`
	Height          int               `json:"height,omitempty"`
	Duration        float64           `json:"duration,omitempty"`
	Bitrate         int               `json:"bitrate,omitempty"`
	Codec           string            `json:"codec,omitempty"`
	Size            int64             `json:"size,omitempty"`
	ContentType     string            `json:"contentType,omitempty"`
	StreamInfo      *StreamInfo       `json:"streamInfo,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	Status          string            `json:"status,omitempty"`
	Thumbnail       *Image            `json:"thumbnail,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type ImageType string

type Image struct {
	ID              string            `json:"id"`
	FileName        string            `json:"fileName"`
	URL             string            `json:"url"`
	Type            ImageType         `json:"type"`
	StorageLocation *S3Object         `json:"storageLocation,omitempty"`
	Width           int               `json:"width,omitempty"`
	Height          int               `json:"height,omitempty"`
	Size            int64             `json:"size,omitempty"`
	ContentType     string            `json:"contentType,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type S3Object struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type StreamInfo struct {
	DownloadURL *string `json:"downloadUrl,omitempty"`
	CdnPrefix   *string `json:"cdnPrefix,omitempty"`
	PlayURL     *string `json:"playUrl,omitempty"`
}

type Credit struct {
	Role     string `json:"role"`
	Name     string `json:"name"`
	PersonID string `json:"personId,omitempty"`
}

type PublishRule struct {
	PublishAt   time.Time `json:"publishAt,omitempty"`
	UnpublishAt time.Time `json:"unpublishAt,omitempty"`
	Regions     []string  `json:"regions,omitempty"`
	AgeRating   string    `json:"ageRating,omitempty"`
}

func (a *Asset) Status() string {
	if a.PublishRule == nil {
		return "draft"
	}

	now := time.Now().UTC()

	if a.PublishRule.PublishAt.IsZero() {
		return "draft"
	}

	if now.Before(a.PublishRule.PublishAt) {
		return "scheduled"
	}

	if !a.PublishRule.UnpublishAt.IsZero() && now.After(a.PublishRule.UnpublishAt) {
		return "expired"
	}

	return "published"
}
