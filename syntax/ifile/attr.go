package ifile

import (
	"path"
	"path/filepath"
	"strings"
)

// Dir get dir path, without last name.
func Dir(fpath string) string {
	return filepath.Dir(fpath)
}

// Name get file/dir name
func Name(fpath string) string {
	return filepath.Base(fpath)
}

// 获取文件路径中的文件名（包括扩展名）
func Name2(filePath string) string {
	_, fileName := filepath.Split(filePath)
	return fileName
}

// 获取文件路径中的文件名（不包括扩展名）
func NameNoExt(filePath string) string {
	fileName := Name(filePath)
	ext := filepath.Ext(fileName)
	if ext != "" {
		fileName = strings.TrimSuffix(fileName, ext)
	}
	return fileName
}

// FileExt get filename ext. alias of path.Ext()
func FileExt(fpath string) string {
	return path.Ext(fpath)
}

// Suffix get filename ext. alias of path.Ext()
func Suffix(fpath string) string {
	return path.Ext(fpath)
}
