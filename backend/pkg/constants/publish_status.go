package constants

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
