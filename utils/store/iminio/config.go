package iminio

import "time"

// config options
type config struct {
	Endpoint        string
	AccesskeyId     string
	SecretaccessKey string

	UseSSL      bool
	ProxySocks5 string

	Debug bool //  打印日志
	Timeout time.Duration // 超时时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		UseSSL: false,
	}
}
