package base

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func Example() {
	msg := "Hello, 世界"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(string(decoded))
	// Output:
	// SGVsbG8sIOS4lueVjA==
	// Hello, 世界
}

func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64EncodeSafe(src []byte) string {
	z := base64.StdEncoding.EncodeToString(src)
	z = strings.ReplaceAll(z, "+", "-")
	z = strings.ReplaceAll(z, "/", "_")
	z = strings.ReplaceAll(z, "=", "")
	return z
}

func Base64Decode(s string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	} else {
		return decoded, nil
	}
}

// Base64UrlDecode base64 url 的 decode
func Base64UrlDecode(s string) ([]byte, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	} else {
		return decoded, nil
	}
}
