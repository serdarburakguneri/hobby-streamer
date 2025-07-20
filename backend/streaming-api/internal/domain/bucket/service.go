package bucket

import (
	"errors"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
)

type BucketOrganizationService struct{}

func NewBucketOrganizationService() *BucketOrganizationService {
	return &BucketOrganizationService{}
}

func (s *BucketOrganizationService) GetBucketStats(bucket *Bucket) *BucketStats {
	if bucket == nil {
		return nil
	}

	stats := &BucketStats{
		TotalAssets:     bucket.AssetCount(),
		PublicAssets:    len(bucket.GetPublicAssets()),
		ReadyAssets:     len(bucket.GetReadyAssets()),
		AssetsWithVideo: len(bucket.GetAssetsWithVideos()),
		AssetsWithImage: len(bucket.GetAssetsWithImages()),
	}

	return stats
}

func (s *BucketOrganizationService) GetBucketRecommendations(bucket *Bucket, limit int) []*asset.Asset {
	if bucket == nil || limit <= 0 {
		return nil
	}

	recommendations := make([]*asset.Asset, 0, limit)

	if bucket.IsFeatured() {
		recommendations = append(recommendations, s.getFeaturedAssets(limit)...)
	}

	if bucket.IsTrending() {
		recommendations = append(recommendations, s.getTrendingAssets(limit)...)
	}

	if bucket.IsCategory() {
		recommendations = append(recommendations, s.getCategoryAssets(bucket, limit)...)
	}

	return recommendations[:min(len(recommendations), limit)]
}

func (s *BucketOrganizationService) getFeaturedAssets(limit int) []*asset.Asset {
	return nil
}

func (s *BucketOrganizationService) getTrendingAssets(limit int) []*asset.Asset {
	return nil
}

func (s *BucketOrganizationService) getCategoryAssets(bucket *Bucket, limit int) []*asset.Asset {
	return nil
}

type BucketDiscoveryService struct{}

func NewBucketDiscoveryService() *BucketDiscoveryService {
	return &BucketDiscoveryService{}
}

func (s *BucketDiscoveryService) GetBucketsByType(bucketType BucketType, limit int) ([]*Bucket, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	buckets := make([]*Bucket, 0, limit)

	return buckets, nil
}

func (s *BucketDiscoveryService) GetBucketsByAssetType(assetType asset.AssetType, limit int) ([]*Bucket, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	buckets := make([]*Bucket, 0, limit)

	return buckets, nil
}

func (s *BucketDiscoveryService) GetBucketsByGenre(genre asset.Genre, limit int) ([]*Bucket, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	buckets := make([]*Bucket, 0, limit)

	return buckets, nil
}

func (s *BucketDiscoveryService) SearchBuckets(query string, filters *BucketSearchFilters) ([]*Bucket, error) {
	if query == "" && filters == nil {
		return nil, ErrInvalidSearchQuery
	}

	buckets := make([]*Bucket, 0)

	if query != "" {
		buckets = append(buckets, s.searchByQuery(query)...)
	}

	if filters != nil {
		buckets = s.filterResults(buckets, filters)
	}

	return buckets, nil
}

func (s *BucketDiscoveryService) searchByQuery(query string) []*Bucket {
	return nil
}

func (s *BucketDiscoveryService) filterResults(buckets []*Bucket, filters *BucketSearchFilters) []*Bucket {
	filtered := make([]*Bucket, 0)

	for _, bucket := range buckets {
		if s.matchesFilters(bucket, filters) {
			filtered = append(filtered, bucket)
		}
	}

	return filtered
}

func (s *BucketDiscoveryService) matchesFilters(bucket *Bucket, filters *BucketSearchFilters) bool {
	if filters == nil {
		return true
	}

	if filters.BucketType != nil && !bucket.Type().Equals(*filters.BucketType) {
		return false
	}

	if filters.Status != nil && (bucket.Status() == nil || !bucket.Status().Equals(*filters.Status)) {
		return false
	}

	if filters.OnlyActive && !bucket.IsActive() {
		return false
	}

	if filters.MinAssetCount > 0 && bucket.AssetCount() < filters.MinAssetCount {
		return false
	}

	if filters.MaxAssetCount > 0 && bucket.AssetCount() > filters.MaxAssetCount {
		return false
	}

	return true
}

type BucketStats struct {
	TotalAssets     int
	PublicAssets    int
	ReadyAssets     int
	AssetsWithVideo int
	AssetsWithImage int
}

type BucketSearchFilters struct {
	BucketType    *BucketType
	Status        *BucketStatus
	OnlyActive    bool
	MinAssetCount int
	MaxAssetCount int
}

var (
	ErrInvalidBucket      = errors.New("invalid bucket")
	ErrInvalidLimit       = errors.New("invalid limit")
	ErrInvalidSearchQuery = errors.New("invalid search query")
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
