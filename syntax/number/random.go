package random

import (
	crand "crypto/rand"
	"io"
	"math/rand"
	"time"
)


// RandInt generate random int between min and max, maybe min,  not be max
func RandInt(min, max int) int {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

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
