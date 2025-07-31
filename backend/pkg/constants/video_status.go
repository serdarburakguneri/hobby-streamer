package constants

const (
	VideoStatusPending     = "pending"
	VideoStatusAnalyzing   = "analyzing"
	VideoStatusTranscoding = "transcoding"
	VideoStatusReady       = "ready"
	VideoStatusFailed      = "failed"
)

var AllowedVideoStatuses = map[string]struct{}{
	VideoStatusPending:     {},
	VideoStatusAnalyzing:   {},
	VideoStatusTranscoding: {},
	VideoStatusReady:       {},
	VideoStatusFailed:      {},
}

func IsValidVideoStatus(s string) bool {
	_, ok := AllowedVideoStatuses[s]
	return ok
}
