package random

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type letter string

const (
	LetterAbc            = letter("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	LetterAbcLower       = letter("abcdefghijklmnopqrstuvwxyz")
	LetterAbcUpper       = letter("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	LetterNumAndLowAbc   = letter("abcdefghijklmnopqrstuvwxyz0123456789")
	LetterNumAndUpperAbc = letter("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	LetterAll            = letter("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// RandString generate random string
// see https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandString(length int, letter letter) string {
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letter[r.Int63()%int64(len(letter))]
	}
	return string(b)
}

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

// UUIdV4 generate a random UUID of version 4 according to RFC 4122
func UUIdV4() (string, error) {
	uuid := make([]byte, 16)

	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
