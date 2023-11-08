package base

import (
	"encoding/base64"
	"fmt"
	"log"
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

func Base64DecodeAll(b64String string) ([]byte, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(b64String)
	if err != nil {
		decodedBytes, err = base64.URLEncoding.DecodeString(b64String)
		if err != nil {
			decodedBytes, err = base64.StdEncoding.DecodeString(b64String)
			if err != nil {
				log.Println("decode base64 fail:", err.Error())
				return []byte{}, err
			}
		}
	}

	return decodedBytes, nil
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

// ParseB64String 解析 base64 是否正常
func ParseB64String(b64String string, padding bool) string {
	lengthBase64 := len(b64String)
	missingPadding := lengthBase64 % 4
	if missingPadding != 0 {
		if padding {
			b64String = b64String + strings.Repeat("=", missingPadding)
		} else {
			b64String = b64String[:lengthBase64-missingPadding]
		}
	}
	return b64String
}
