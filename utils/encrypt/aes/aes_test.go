package aes

import (
	"log"
	"testing"
)

var key = []byte("secret key 123xx")
var msg = []byte("my message")

// go test -v -run TestEncrypt
func TestEncrypt(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	t.Run("CBC", func(t *testing.T) {
		d, e := EncryptCBC(msg, key, PaddingPkcs7)
		log.Println("EncryptCBC", d, e)
		dd, err := DecryptCBC(d, key, PaddingPkcs7)
		log.Println("DecryptCBC", string(dd), err)
	})

	t.Run("CFB", func(t *testing.T) {
		d, e := EncryptCFB(msg, key)
		log.Println("EncryptCFB", d, e)
		dd, err := DecryptCFB(d, key)
		log.Println("DecryptCBC", string(dd), err)
	})
}
