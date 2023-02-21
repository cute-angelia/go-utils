package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// Deprecated: hash.Hash(hash.AlgoMD5, "txt")
func NewEncodeMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
