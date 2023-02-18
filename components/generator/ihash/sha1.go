package ihash

import (
	"crypto/sha1"
	"encoding/hex"
)

// Deprecated: ihash.Hash(ihash.AlgoSha1, "txt")
func NewEncodeSha1(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
