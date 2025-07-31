package constants

const (
	AssetStatusDraft     = "draft"
	AssetStatusPublished = "published"
	AssetStatusScheduled = "scheduled"
	AssetStatusExpired   = "expired"
)

var AllowedAssetStatuses = map[string]struct{}{
	AssetStatusDraft:     {},
	AssetStatusPublished: {},
	AssetStatusScheduled: {},
	AssetStatusExpired:   {},
}

func IsValidAssetStatus(s string) bool {
	_, ok := AllowedAssetStatuses[s]
	return ok
}
