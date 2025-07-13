package asset

import "github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"

const (
	ErrIDShouldNotBeSet = "id should not be set by client"
	ErrIDMismatch       = "id mismatch"
	ErrImageExists      = "image with same filename already exists"

	AssetTypeMovie           = "movie"
	AssetTypeSeries          = "series"
	AssetTypeSeason          = "season"
	AssetTypeEpisode         = "episode"
	AssetTypeDocumentary     = "documentary"
	AssetTypeMusic           = "music"
	AssetTypePodcast         = "podcast"
	AssetTypeTrailer         = "trailer"
	AssetTypeBehindTheScenes = "behind_the_scenes"
	AssetTypeInterview       = "interview"

	AssetGenreAction      = "action"
	AssetGenreDrama       = "drama"
	AssetGenreComedy      = "comedy"
	AssetGenreHorror      = "horror"
	AssetGenreSciFi       = "sci_fi"
	AssetGenreRomance     = "romance"
	AssetGenreThriller    = "thriller"
	AssetGenreFantasy     = "fantasy"
	AssetGenreDocumentary = "documentary"
	AssetGenreMusic       = "music"
	AssetGenreNews        = "news"
	AssetGenreSports      = "sports"
	AssetGenreKids        = "kids"
	AssetGenreEducational = "educational"
)

const (
	VideoStatusPending          = "pending"
	VideoStatusAnalyzing        = "analyzing"
	VideoStatusTranscoding      = "transcoding"
	VideoStatusReady            = "ready"
	VideoStatusFailed           = "failed"
	VideoStatusAnalyzeCompleted = "analyze_completed"
	VideoStatusAnalyzeFailed    = "analyze_failed"
)

const (
	VideoTypeMain      VideoType = "main"
	VideoTypeTrailer   VideoType = "trailer"
	VideoTypeBehind    VideoType = "behind_the_scenes"
	VideoTypeInterview VideoType = "interview"
)

const (
	VideoVariantRaw  = "raw"
	VideoVariantHLS  = "hls"
	VideoVariantDASH = "dash"
)

var (
	AssetStatusDraft     = constants.AssetStatusDraft
	AssetStatusPublished = constants.AssetStatusPublished
)
