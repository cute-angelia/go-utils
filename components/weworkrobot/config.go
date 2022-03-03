package weworkrobot

import "os"

// PackageName ..
const PackageName = "contrib.qcloud.weworkrobot"

// config
type config struct {
	Uri                 string
	Key                 string
	MentionedList       []string
	MentionedMobileList []string
	Debug               bool

	From     string // 来源
	Topic    string // 主题
	WithTime bool   // 文本带时间
	Retry    int    // 重试次数
}

// DefaultConfig ...
func DefaultConfig() *config {
	hostname, _ := os.Hostname()
	c := config{
		From:     hostname,
		Uri:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send",
		Key:      "default",
		WithTime: true,
		Retry:    3,
	}
	return &c
}
