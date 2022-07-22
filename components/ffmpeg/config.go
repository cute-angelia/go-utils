package ffmpeg

import "time"

var PackageName = "component.utils.ffmpeg"

// config options
type config struct {
	FfmpegPath string
	Timeout    time.Duration
	FilesPath  string // 源文件路径
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		FfmpegPath: "/usr/local/bin/ffmpeg",
		Timeout:    time.Second * 60,
	}
}
