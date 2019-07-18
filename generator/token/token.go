package token

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func NewToken(secret string) *Token {

	if len(secret) == 0 {
		secret = "d923lxll32920392093232323"
	}

	return &Token{
		Secret: secret,
	}
}

type Token struct {
	Secret string
}

func (s Token) Generate(uid int, timestamp int) string {
	str := fmt.Sprintf("%d_%d_%s", uid, timestamp, s.Secret)
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func (s Token) Check(token string, uid int, timestamp int) bool {
	ztoken := s.Generate(uid, timestamp)
	if token == ztoken {
		return true
	} else {
		return false
	}
}
