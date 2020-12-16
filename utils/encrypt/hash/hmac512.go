package hash

import (
	"crypto/sha512"
	"crypto/hmac"
	"encoding/base64"
)

func NewEncodeHmac512(text string, secret string) string {
	hmac512 := hmac.New(sha512.New, []byte(secret))
	hmac512.Write([]byte(text))

	return base64.StdEncoding.EncodeToString(hmac512.Sum(nil))
}
