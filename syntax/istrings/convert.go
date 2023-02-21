package istrings

import "unsafe"

// 对比 byte 与 字符串
func UnsafeEqual(a string, b []byte) bool {
	bbp := *(*string)(unsafe.Pointer(&b))
	return a == bbp
}

// byte[] 转 string
func UnsafeConvertString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
