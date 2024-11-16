package random

import (
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
