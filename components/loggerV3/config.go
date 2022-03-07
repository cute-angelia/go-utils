package loggerV3

import "github.com/rs/zerolog"

const ComponentName = "Component.LoggerV3"

type config struct {
	project    string
	isOnline   bool // online ture is file , false is stdout, 记录文件或者命令行输出
	fileJson   bool // 记录文件的输出格式：默认 json 或者切换为 命令行输出格式
	maxSize    int
	maxBackups int  // log nums
	maxAge     int  // days
	everyday   bool // log every day
	level      zerolog.Level
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		isOnline:   true,
		maxSize:    100,
		maxBackups: 10,
		maxAge:     10,
		everyday:   true,
		level:      zerolog.DebugLevel,
		fileJson:   true,
		project:    "default",
	}
}
