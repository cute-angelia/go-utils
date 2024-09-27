package istrings

import (
	"fmt"
	"strconv"
	"unsafe"
)

var ConvertUtil = convertUtil{}

// convertUtil 转换工具
type convertUtil struct{}

// UnsafeEqual 对比 byte 与 字符串
func (that convertUtil) UnsafeEqual(a string, b []byte) bool {
	bbp := *(*string)(unsafe.Pointer(&b))
	return a == bbp
}

// UnsafeConvertString  byte[] 转 string
func (that convertUtil) UnsafeConvertString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Int 转int
func String2Int[T int | int32 | int64](strint string) (T, error) {
	if strint == "" {
		return 0, fmt.Errorf("str is empty")
	}
	uid, err := strconv.ParseInt(strint, 10, 64)
	if err != nil {
		return 0, err
	}
	switch any(T(0)).(type) {
	case int:
		return T(int(uid)), nil
	case int32:
		return T(int32(uid)), nil
	case int64:
		return T(uid), nil
	default:
		return 0, fmt.Errorf("unsupport %v", any(T(0)))
	}
}
