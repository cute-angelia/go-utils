package hash

import (
	"encoding/hex"
	"crypto/sha256"
)

func NewEncodeSha256(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
