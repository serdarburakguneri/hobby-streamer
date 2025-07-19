package asset

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func generateID() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000_000))
	if err != nil {
		return "000000000"
	}
	return strconv.Itoa(int(n.Int64()))
}

func getFileExtension(format string) string {
	switch format {
	case "hls":
		return "m3u8"
	case "dash":
		return "mpd"
	default:
		return "mp4"
	}
}
