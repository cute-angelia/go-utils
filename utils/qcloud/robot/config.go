package robot

// PackageName ..
const PackageName = "contrib.qcloud.robot"

// config
type config struct {
	Uri                 string
	Key                 string
	MentionedList       []string
	MentionedMobileList []string
	Debug bool
}

// DefaultConfig ...
func DefaultConfig() *config {
	c := config{
		Uri: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send",
		Key: "default",
	}
	return &c
}
