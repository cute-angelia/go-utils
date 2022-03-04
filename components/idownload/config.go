package idownload

import (
	"runtime"
	"time"
)

// config options
type config struct {
	ProxyHttp                string // 代理 http://ip:port
	ProxySocks5              string // 代理 ip:port
	Cookie                   string // cookie
	Referer                  string // Referer
	UserAgent                string // user-agent
	UseRandomUserAgent       bool
	UseRandomUserAgentMobile bool
	Authorization            string

	// 并发数量
	Concurrency int
	// 分片下载
	Resume bool

	Timeout time.Duration // 超时时间
	Debug   bool          //  debug 日志
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	concurrency := runtime.NumCPU()
	return &config{
		Concurrency:              concurrency,
		Resume:                   true,
		ProxyHttp:                "",
		ProxySocks5:              "",
		Cookie:                   "",
		UserAgent:                "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		UseRandomUserAgent:       false,
		UseRandomUserAgentMobile: false,
		Referer:                  "",
		Debug:                    false,
		Authorization:            "",
	}
}