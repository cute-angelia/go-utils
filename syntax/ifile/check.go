package ifile

import (
	"bytes"
	"os"
	"path"
)

var (
	// perm and flags for create log file
	DefaultDirPerm  os.FileMode = 0775
	DefaultFilePerm os.FileMode = 0664

	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

// alias methods
var (
	DirExist  = IsDir
	FileExist = IsFile
	PathExist = PathExists
)

// PathExists reports whether the named file or directory exists.
func PathExists(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsDir reports whether the named directory exists.
func IsDir(path string) bool {
	if path == "" {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

// path exist
func IsExistPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// FileExists reports whether the named file or directory exists.
func IsExistFile(path string) bool {
	return IsFile(path)
}

// IsFile reports whether the named file or directory exists.
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
	"bmp": "image/bmp",
	"gif": "image/gif",
	"ief": "image/ief",
	"jpg": "image/jpeg",
	// "jpe":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"webp": "image/webp",
}

// IsImageFile check file is image file.
func IsImageFile(path string) bool {
	mime := MimeType(path)
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

// IsZipFile check is zip file.
// from https://blog.csdn.net/wangshubo1989/article/details/71743374
func IsZipFile(filepath string) bool {
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
