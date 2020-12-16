package hash

import (
	"encoding/hex"
	"crypto/sha512"
)

func NewEncodeSha512(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
