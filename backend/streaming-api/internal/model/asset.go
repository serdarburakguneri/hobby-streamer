package model

import "time"

type Asset struct {
	ID          string       `json:"id"`
	Slug        string       `json:"slug"`
	Title       *string      `json:"title,omitempty"`
	Description *string      `json:"description,omitempty"`
	Type        string       `json:"type"`
	Genre       *string      `json:"genre,omitempty"`
	Genres      []string     `json:"genres,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Status      *string      `json:"status,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Metadata    *string      `json:"metadata,omitempty"`
	OwnerID     *string      `json:"ownerId,omitempty"`
	Videos      []Video      `json:"videos,omitempty"`
	Images      []Image      `json:"images,omitempty"`
	PublishRule *PublishRule `json:"publishRule,omitempty"`
}

type Video struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	Format          string      `json:"format"`
	StorageLocation S3Object    `json:"storageLocation"`
	Width           *int        `json:"width,omitempty"`
	Height          *int        `json:"height,omitempty"`
	Duration        *float64    `json:"duration,omitempty"`
	Bitrate         *int        `json:"bitrate,omitempty"`
	Codec           *string     `json:"codec,omitempty"`
	Size            *int        `json:"size,omitempty"`
	ContentType     *string     `json:"contentType,omitempty"`
	StreamInfo      *StreamInfo `json:"streamInfo,omitempty"`
	Metadata        *string     `json:"metadata,omitempty"`
	Status          *string     `json:"status,omitempty"`
	Thumbnail       *Image      `json:"thumbnail,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

type Image struct {
	ID              string      `json:"id"`
	FileName        string      `json:"fileName"`
	URL             string      `json:"url"`
	Type            string      `json:"type"`
	StorageLocation *S3Object   `json:"storageLocation,omitempty"`
	Width           *int        `json:"width,omitempty"`
	Height          *int        `json:"height,omitempty"`
	Size            *int        `json:"size,omitempty"`
	ContentType     *string     `json:"contentType,omitempty"`
	StreamInfo      *StreamInfo `json:"streamInfo,omitempty"`
	Metadata        *string     `json:"metadata,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

type S3Object struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type StreamInfo struct {
	DownloadURL *string `json:"downloadUrl,omitempty"`
	CDNPrefix   *string `json:"cdnPrefix,omitempty"`
	URL         *string `json:"url,omitempty"`
}

type PublishRule struct {
	IsPublic    bool       `json:"isPublic"`
	PublishAt   *time.Time `json:"publishAt,omitempty"`
	UnpublishAt *time.Time `json:"unpublishAt,omitempty"`
	Regions     []string   `json:"regions,omitempty"`
	AgeRating   *string    `json:"ageRating,omitempty"`
}

type AssetPage struct {
	Items   []Asset `json:"items"`
	NextKey *string `json:"nextKey,omitempty"`
}
