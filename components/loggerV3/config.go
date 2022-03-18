package loggerV3

import "github.com/rs/zerolog"

const ComponentName = "Component.LoggerV3"

type config struct {
	Project    string `json:"project"`
	IsOnline   bool   // online ture is file , false is stdout, 记录文件或者命令行输出
	FileJson   bool   // 记录文件的输出格式：默认 json 或者切换为 命令行输出格式
	MaxSize    int
	MaxBackups int  // log nums
	MaxAge     int  // days
	Everyday   bool // log every day
	Level      zerolog.Level
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Project:    "default",
		IsOnline:   true,
		FileJson:   true,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     10,
		Everyday:   true,
		Level:      zerolog.DebugLevel,
	}
}
