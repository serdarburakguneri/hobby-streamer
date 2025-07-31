package operations

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
