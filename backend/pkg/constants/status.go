package constants

const (
	StatusPending = "pending"
	StatusActive  = "active"
	StatusFailed  = "failed"

	AssetStatusDraft     = "draft"
	AssetStatusPublished = "published"
	AssetStatusScheduled = "scheduled"
	AssetStatusExpired   = "expired"

	VideoStatusPending     = "pending"
	VideoStatusAnalyzing   = "analyzing"
	VideoStatusTranscoding = "transcoding"
	VideoStatusReady       = "ready"
	VideoStatusFailed      = "failed"
)

const (
	AgeRatingG    = "G"
	AgeRatingPG   = "PG"
	AgeRatingPG13 = "PG-13"
	AgeRatingR    = "R"
	AgeRatingNC17 = "NC-17"
	AgeRatingTVY  = "TV-Y"
	AgeRatingTVY7 = "TV-Y7"
	AgeRatingTVG  = "TV-G"
	AgeRatingTVPG = "TV-PG"
	AgeRatingTV14 = "TV-14"
	AgeRatingTVMA = "TV-MA"
)

const (
	VideoQualityMain = "main"
	VideoQualitySD   = "sd"
	VideoQualityHD   = "hd"
	VideoQualityFHD  = "fhd"
	VideoQuality4K   = "4k"
)

const (
	VideoTypeMain      = "main"
	VideoTypeTrailer   = "trailer"
	VideoTypeBehind    = "behind"
	VideoTypeInterview = "interview"
)

const (
	VideoStreamingFormatRaw  = "raw"
	VideoStreamingFormatHLS  = "hls"
	VideoStreamingFormatDASH = "dash"
)



type PublishStatus int

const (
	PublishStatusInvalid PublishStatus = iota
	PublishStatusNotReady
	PublishStatusNotConfigured
	PublishStatusScheduled
	PublishStatusPublished
	PublishStatusExpired
	PublishStatusDraft
)
