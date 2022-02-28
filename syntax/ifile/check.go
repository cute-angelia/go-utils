package ifile

import (
	"bytes"
	"os"
	"path"
)

var (
	// DefaultDirPerm perm and flags for create log file
	DefaultDirPerm  os.FileMode = 0775
	DefaultFilePerm os.FileMode = 0664
)

// IsExist reports whether the named file or directory exists.
func IsExist(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// IsDir reports whether the named file or directory exists.
func IsDir(path string) bool {
	if path == "" {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

// IsFile reports whether the named file exists.
func IsFile(path string) bool {
	if path == "" {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return !fi.IsDir()
	}
	return false
}

// IsAbsPath is abs path.
func IsAbsPath(aPath string) bool {
	return path.IsAbs(aPath)
}

// ImageMimeTypes refer net/http package
var ImageMimeTypes = map[string]string{
	"bmp":  "image/bmp",
	"gif":  "image/gif",
	"ief":  "image/ief",
	"jpg":  "image/jpeg",
	"jpe":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"webp": "image/webp",
}

const (
	CheckTypeExt = iota
	CheckTypeMimeType
)

// IsImage is checked file is Image
func IsImage(uri string, checkType int32) bool {
	if checkType == CheckTypeExt {
		uri = NewFileName(uri).CleanUrl()
		ext := path.Ext(uri)
		if ext == ".jpg" ||
			ext == ".png" ||
			ext == ".svg" ||
			ext == ".gif" ||
			ext == ".jpeg" ||
			ext == ".webp" ||
			ext == ".icon" {
			return true
		} else {
			return false
		}
	}

	if checkType == CheckTypeMimeType {
		mime := MimeType(uri)
		if mime == "" {
			return false
		}

		for _, imgMime := range ImageMimeTypes {
			if imgMime == mime {
				return true
			}
		}
		return false
	}
	return false
}

// IsZip check is zip file.
// from https://blog.csdn.net/wangshubo1989/article/details/71743374
func IsZip(filepath string) bool {
	f, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

// IsCanWrite 测试写权限
func IsCanWrite(inFile string) bool {
	// 这个例子测试写权限，如果没有写权限则返回error。
	// 注意文件不存在也会返回error，需要检查error的信息来获取到底是哪个错误导致。
	file, err := os.OpenFile(inFile, os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
	}
	file.Close()
	return true
}

// IsCanRead 测试读权限
func IsCanRead(inFile string) bool {
	file, err := os.OpenFile(inFile, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
	}
	file.Close()
	return true
}
