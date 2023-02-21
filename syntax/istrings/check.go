package istrings

// https://github.com/gookit/goutil/blob/master/strutil/check.go

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	Equal     = strings.EqualFold
	HasPrefix = strings.HasPrefix
	HasSuffix = strings.HasSuffix
)

// refer from github.com/yuin/goldmark/util
var spaceTable = [256]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// 检查字符串是否是数字字符串
func IsNumberString(str string) (int64, error) {
	if v, err := strconv.ParseInt(str, 10, 64); err == nil {
		// fmt.Printf("%q looks like a number.\n", v)
		return v, nil
	} else {
		return 0, err
	}
}

// IsNumeric returns true if the given character is a numeric, otherwise false.
func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsAlphabet char
func IsAlphabet(char uint8) bool {
	// A 65 -> Z 90
	if char >= 'A' && char <= 'Z' {
		return true
	}

	// a 97 -> z 122
	if char >= 'a' && char <= 'z' {
		return true
	}

	return false
}

// IsAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func IsAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// StrPos alias of the strings.Index
func StrPos(s, sub string) int {
	return strings.Index(s, sub)
}

// BytePos alias of the strings.IndexByte
func BytePos(s string, bt byte) int {
	return strings.IndexByte(s, bt)
}

// RunePos alias of the strings.IndexRune
func RunePos(s string, ru rune) int {
	return strings.IndexRune(s, ru)
}

// IsStartOf alias of the strings.HasPrefix
func IsStartOf(s, sub string) bool {
	return strings.HasPrefix(s, sub)
}

// IsEndOf alias of the strings.HasSuffix
func IsEndOf(s, sub string) bool {
	return strings.HasSuffix(s, sub)
}

// Len of the string
func Len(s string) int {
	return len(s)
}

// Utf8Len of the string
func Utf8len(s string) int {
	return utf8.RuneCount([]byte(s))
}

// ValidUtf8String check
func ValidUtf8String(s string) bool {
	return utf8.ValidString(s)
}

// ----- refer from github.com/yuin/goldmark/util

// IsSpace returns true if the given character is a space, otherwise false.
func IsSpace(c byte) bool {
	return spaceTable[c] == 1
}

// IsSpaceRune returns true if the given rune is a space, otherwise false.
func IsSpaceRune(r rune) bool {
	return r <= 256 && IsSpace(byte(r)) || unicode.IsSpace(r)
}

// IsBlank returns true if the given string is all space characters.
func IsBlank(s string) bool {
	return IsBlankBytes([]byte(s))
}

// IsBlankBytes returns true if the given []byte is all space characters.
func IsBlankBytes(bs []byte) bool {
	for _, b := range bs {
		if !IsSpace(b) {
			return false
		}
	}
	return true
}

// IsSymbol reports whether the rune is a symbolic character.
func IsSymbol(r rune) bool {
	return unicode.IsSymbol(r)
}
