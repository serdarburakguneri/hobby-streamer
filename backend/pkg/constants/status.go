package constants

const (
	StatusPending = "pending"
	StatusActive  = "active"
	StatusFailed  = "failed"
)

var AllowedStatuses = map[string]struct{}{
	StatusPending: {},
	StatusActive:  {},
	StatusFailed:  {},
}

func IsValidStatus(s string) bool {
	_, ok := AllowedStatuses[s]
	return ok
}
