package apicache

import "time"

const PackageName = "component.http.apicache"

// config options
type config struct {
	Timeout  time.Duration // 失效时间
	DbName   string
	CacheKey string

	OnlyToday bool // 凌晨刷新
	Debug     bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Timeout:   time.Minute * 10,
		DbName:    "cache",
		OnlyToday: false,
		Debug:     false,
	}
}
