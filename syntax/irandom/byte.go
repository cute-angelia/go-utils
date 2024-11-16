package irandom

import (
	crand "crypto/rand"
	"io"
)

// RandBytes generate random byte slice
func RandBytes(length int) []byte {
	if length < 1 {
		return []byte{}
	}
	b := make([]byte, length)

	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return nil
	}
	return b
}
