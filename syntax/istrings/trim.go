package istrings

import (
	"strings"
	"unicode"
)

func TrimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}