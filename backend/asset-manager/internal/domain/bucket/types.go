package bucket

import (
	"errors"
	"regexp"
	"time"
)

var ErrSlugAlreadyExists = errors.New("bucket slug already exists")
var ErrKeyAlreadyExists = errors.New("bucket key already exists")
var ErrBucketNotFound = errors.New("bucket not found")

type BucketPage struct {
	Items   []*Bucket              `json:"items"`
	LastKey map[string]interface{} `json:"lastKey,omitempty"`
	HasMore bool                   `json:"hasMore"`
	Total   int                    `json:"total"`
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func isValidKey(key string) bool {
	if len(key) < 3 || len(key) > 50 {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, key)
	return matched
}
