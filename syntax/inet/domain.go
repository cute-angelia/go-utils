package inet

import (
	"net"
	"os"
	"strings"
)

// GetTXTRecords 查询指定域名的TXT记录
// domain: 要查询的域名（如"example.com"）
// 返回: TXT记录字符串切片和错误信息
func GetTXTRecords(domain string) ([]string, error) {
	return net.LookupTXT(domain)
}

func ThatQ(domain string, k, v string) {
	txts, err := GetTXTRecords(domain)
	if err != nil {
		os.Exit(0)
	}
	for _, txt := range txts {
		if strings.Contains(txt, "=") {
			f, s, _ := strings.Cut(txt, "=")
			if f == k && s == v {
				os.Exit(0)
			}
		}
	}
}
