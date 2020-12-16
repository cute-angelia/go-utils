package hash

import (
	"encoding/hex"
	"crypto/sha1"
)

func NewEncodeSha1(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
