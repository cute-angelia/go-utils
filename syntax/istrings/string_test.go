package istrings

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	str := ".minio.sys/buckets/social/comic/藏锋行/102_化身修罗/藏锋行_化身修罗_1629888559534437000.jpg"
	t.Log(strings.Contains(str, ".minio.sys"))
}
