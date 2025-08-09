package constants

const (
	BucketTypeCollection = "collection"
	BucketTypePlaylist   = "playlist"
	BucketTypeCategory   = "category"
	BucketTypeFeatured   = "featured"
	BucketTypeTrending   = "trending"
)

var AllowedBucketTypes = map[string]struct{}{
	BucketTypeCollection: {},
	BucketTypePlaylist:   {},
	BucketTypeCategory:   {},
	BucketTypeFeatured:   {},
	BucketTypeTrending:   {},
}

func IsValidBucketType(t string) bool {
	_, ok := AllowedBucketTypes[t]
	return ok
}
