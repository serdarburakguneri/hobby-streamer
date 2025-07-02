package bucket

import (
	"time"
)

type Bucket struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Name        string    `json:"name" dynamodbav:"name"`
	Description string    `json:"description,omitempty" dynamodbav:"description"`
	AssetIDs    []string  `json:"assetIds,omitempty" dynamodbav:"asset_ids"`
	CreatedAt   time.Time `json:"createdAt" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" dynamodbav:"updated_at"`
}
