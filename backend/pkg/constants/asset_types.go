package constants

const (
	AssetTypeMovie        = "movie"
	AssetTypeTVShow       = "tv_show"
	AssetTypeSeries       = "series"
	AssetTypeSeason       = "season"
	AssetTypeEpisode      = "episode"
	AssetTypeDocumentary  = "documentary"
	AssetTypeShort        = "short"
	AssetTypeTrailer      = "trailer"
	AssetTypeBonus        = "bonus"
	AssetTypeBehindScenes = "behind_scenes"
	AssetTypeInterview    = "interview"
	AssetTypeMusicVideo   = "music_video"
	AssetTypePodcast      = "podcast"
	AssetTypeLive         = "live"
)

var AllowedAssetTypes = map[string]struct{}{
	AssetTypeMovie:        {},
	AssetTypeTVShow:       {},
	AssetTypeSeries:       {},
	AssetTypeSeason:       {},
	AssetTypeEpisode:      {},
	AssetTypeDocumentary:  {},
	AssetTypeShort:        {},
	AssetTypeTrailer:      {},
	AssetTypeBonus:        {},
	AssetTypeBehindScenes: {},
	AssetTypeInterview:    {},
	AssetTypeMusicVideo:   {},
	AssetTypePodcast:      {},
	AssetTypeLive:         {},
}

func IsValidAssetType(t string) bool {
	_, ok := AllowedAssetTypes[t]
	return ok
}
