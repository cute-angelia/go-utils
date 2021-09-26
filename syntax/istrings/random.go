package istrings

import (
	mathRand "math/rand"
	"time"
)

const (
	AlphaBet  = "abcdefghijklmnopqrstuvwxyz"
	AlphaNum  = "abcdefghijklmnopqrstuvwxyz0123456789"
	AlphaNum2 = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandomChars generate give length random chars at `a-z`
func RandomChars(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(25) // 0 - 25
		cs[i] = AlphaBet[idx]
	}

	return string(cs)
}

// RandomCharsV2 generate give length random chars in `0-9a-z`
func RandomCharsV2(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(35) // 0 - 35
		cs[i] = AlphaNum[idx]
	}

	return string(cs)
}

// RandomCharsV3 generate give length random chars in `0-9a-zA-Z`
func RandomCharsV3(ln int) string {
	cs := make([]byte, ln)
	for i := 0; i < ln; i++ {
		// 1607400451937462000
		mathRand.Seed(time.Now().UnixNano())
		idx := mathRand.Intn(61) // 0 - 61
		cs[i] = AlphaNum2[idx]
	}

	return string(cs)
}
