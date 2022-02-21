package islice

// Clone
// 1.18 已支持，这里兼容旧版 不方便升级的包
func Clone(b []byte) []byte {
	return append([]byte(nil), b...)
}

// CloneCopy = Clone
func CloneCopy(b []byte) []byte {
	b2 := make([]byte, len(b))
	copy(b2, b)
	return b2
}
