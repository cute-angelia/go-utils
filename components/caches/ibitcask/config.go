package ibitcask

import "time"

// config options
type config struct {
	Path            string
	MaxDatafileSize int
	GcInterval      time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Path:            "./bitcask",
		MaxDatafileSize: 1024 * 1024 * 200,
		GcInterval:      time.Hour,
	}
}
