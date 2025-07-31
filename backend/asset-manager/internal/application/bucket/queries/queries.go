package queries

import "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"

type GetBucketQuery struct {
	ID valueobjects.BucketID
}

type GetBucketByKeyQuery struct {
	Key valueobjects.BucketKey
}

type ListBucketsQuery struct {
	Limit  *int
	Offset *int
}

type SearchBucketsQuery struct {
	Query  string
	Limit  *int
	Offset *int
}

type GetBucketsByOwnerQuery struct {
	OwnerID valueobjects.OwnerID
	Limit   *int
	Offset  *int
}

type GetBucketAssetsQuery struct {
	BucketID valueobjects.BucketID
	Limit    *int
}
