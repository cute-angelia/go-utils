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

func Base64Encode(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func Base64EncodeSafe(in string) string {
	z := base64.StdEncoding.EncodeToString([]byte(in))
	z = strings.ReplaceAll(z, "+", "-")
	z = strings.ReplaceAll(z, "/", "_")
	z = strings.ReplaceAll(z, "=", "")
	return z
}

func Base64Decode(in string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", err
	} else {
		return string(decoded), nil
	}
}
