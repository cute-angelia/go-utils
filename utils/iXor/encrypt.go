package iXor

import "encoding/base64"

// XorEncrypt 加密
// 返回 base64
func XorEncrypt(input []byte, key string) string {
	var result string
	b := base64.StdEncoding.EncodeToString(input)
	for i := 0; i < len(b); i++ {
		j := i % len(key)
		// 异或运算并转换为字符
		c := string(b[i] ^ key[j])
		result += c
	}
	return base64.StdEncoding.EncodeToString([]byte(result))
}

func XorDecrypt(input string, key string) string {
	xorInput, _ := base64.StdEncoding.DecodeString(input)
	var result string
	for i := 0; i < len(xorInput); i++ {
		j := i % len(key)
		// 异或运算并转换为字符
		c := string(xorInput[i] ^ key[j])
		result += c
	}
	rs, _ := base64.StdEncoding.DecodeString(result)
	return string(rs)
}
