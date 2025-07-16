package model

import "time"

type Bucket struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Type        string    `json:"type"`
	Status      *string   `json:"status,omitempty"`
	AssetIDs    []string  `json:"assetIds,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Assets      []Asset   `json:"assets,omitempty"`
}

type BucketPage struct {
	Items   []Bucket `json:"items"`
	NextKey *string  `json:"nextKey,omitempty"`
}
