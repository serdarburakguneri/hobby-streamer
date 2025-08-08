package events

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func BuildJobCorrelationID(assetID, videoID, jobType, format, quality string) string {
	parts := []string{strings.ToLower(jobType), strings.ToLower(format), strings.ToLower(quality), assetID, videoID}
	key := strings.Join(parts, ":")
	h := sha1.Sum([]byte(key))
	return fmt.Sprintf("job-%s", hex.EncodeToString(h[:8]))
}
