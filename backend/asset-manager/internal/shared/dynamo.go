package shared

import (
	"time"
)

func NowUTCString() string {
	return time.Now().UTC().Format(time.RFC3339)
}
