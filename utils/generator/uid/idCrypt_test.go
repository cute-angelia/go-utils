package uid

import (
	"testing"
	"log"
)

func TestAll(t *testing.T) {
	in32 := uint32(233)
	in64 := uint64(233)

	crypt := NewIdCrypt32()
	log.Println(crypt.Encrypt(in32))
	log.Println(crypt.Decrypt(crypt.Encrypt(in32)))

	crypt2 := NewIdCrypt64()
	log.Println(crypt2.Encrypt(in64))
	log.Println(crypt2.Decrypt(crypt2.Encrypt(in64)))
}