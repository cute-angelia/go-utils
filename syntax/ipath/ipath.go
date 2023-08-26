package ipath

import (
	"path"
	"strings"
)

// Clean cleans up given path and returns a relative path that goes straight down.
func Clean(p string) string {
	return strings.Trim(path.Clean("/"+p), "/")
}

// GetFileName 根据路径获取文件名
func GetFileName(filePath string) (name string, ext string) {
	ext = path.Ext(filePath)
	name = path.Base(filePath[:len(filePath)-len(ext)])
	return
}
