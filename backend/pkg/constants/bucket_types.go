package constants

const (
	BucketTypeCollection = "collection"
	BucketTypePlaylist   = "playlist"
	BucketTypeCategory   = "category"
)

var AllowedBucketTypes = map[string]struct{}{
	BucketTypeCollection: {},
	BucketTypePlaylist:   {},
	BucketTypeCategory:   {},
}

func IsValidBucketType(t string) bool {
	_, ok := AllowedBucketTypes[t]
	return ok
}
