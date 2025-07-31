package valueobjects

import "time"

type Timestamps struct {
	createdAt time.Time
	updatedAt time.Time
}

func NewTimestamps() *Timestamps {
	now := time.Now().UTC()
	return &Timestamps{createdAt: now, updatedAt: now}
}

func (t *Timestamps) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Timestamps) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Timestamps) Update() {
	t.updatedAt = time.Now().UTC()
}

type AssetID = ID

func MustNewAssetID(value string) AssetID {
	id, _ := NewAssetID(value)
	return *id
}

type VideoID = ID
type ImageID = ID

func ParseS3URL(url string) (string, string, error) {
	return parseS3URL(url)
}
