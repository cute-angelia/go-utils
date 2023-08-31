package ipath

import (
	"path"
	"strings"
)

// Clean cleans up given path and returns a relative path that goes straight down.
func Clean(p string) string {
	v, _, _ := strings.Cut(p, "?")
	return strings.Trim(path.Clean("/"+v), "/")
}

// GetFileName 根据路径获取文件名
func GetFileName(filePath string) (name string) {
	filePath = Clean(filePath)
	ext := strings.ToLower(path.Ext(filePath))
	name = path.Base(filePath[:len(filePath)-len(ext)])
	return
}

// GetFileNameAndExt 根据路径获取文件名 和 后缀
func GetFileNameAndExt(filePath string) (name string, ext string) {
	filePath = Clean(filePath)
	ext = strings.ToLower(path.Ext(filePath))
	name = path.Base(filePath[:len(filePath)-len(ext)])
	return
}

// GetFileExt 后缀
func GetFileExt(filePath string) (ext string, found bool) {
	filePath = Clean(filePath)
	ext = strings.ToLower(path.Ext(filePath))
	if strings.Contains(ext, ".") {
		return ext, true
	} else {
		return ext, found
	}
}
