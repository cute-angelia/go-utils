package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Encode(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
