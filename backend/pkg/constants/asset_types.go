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
	VideoFormatHLS        = "hls"
	VideoFormatDASH       = "dash"
	VideoFormatRAW        = "raw"
	VideoFormatMP4        = "mp4"
	VideoFormatWEBM       = "webm"
	VideoFormatAVI        = "avi"
	VideoFormatMOV        = "mov"
	VideoFormatMKV        = "mkv"
)

var AllowedAssetTypes = []string{
	AssetTypeMovie,
	AssetTypeTVShow,
	AssetTypeSeries,
	AssetTypeSeason,
	AssetTypeEpisode,
	AssetTypeDocumentary,
	AssetTypeShort,
	AssetTypeTrailer,
	AssetTypeBonus,
	AssetTypeBehindScenes,
	AssetTypeInterview,
	AssetTypeMusicVideo,
	AssetTypePodcast,
	AssetTypeLive,
}

var AllowedVideoFormats = []string{
	VideoFormatHLS,
	VideoFormatDASH,
	VideoFormatRAW,
	VideoFormatMP4,
	VideoFormatWEBM,
	VideoFormatAVI,
	VideoFormatMOV,
	VideoFormatMKV,
}
