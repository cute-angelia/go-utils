package iXor

import (
	"log"
	"testing"
)

func TestName(t *testing.T) {
	var key = "helloxx.zl"
	var a = "hello 我爱你中国 🤔️"
	var en = XorEncrypt([]byte(a), key)
	var de = XorDecrypt(en, key)
	log.Println(en)
	log.Println(de)
}
