package weworkrobot

import (
	"fmt"
	"github.com/guonaihong/gout"
	"log"
	"time"
)

type Component struct {
	config *config
}

func newComponent(cfg *config) *Component {
	return &Component{
		config: cfg,
	}
}

func (self Component) generateContent(content string) string {
	timenow := ""
	if self.config.WithTime {
		timenow = time.Now().Format("2006-01-02 15:04:05")
	}

	topic := ""
	if len(self.config.Topic) > 0 {
		topic = " [" + self.config.Topic + "]"
	}

	from := ""
	if len(self.config.From) > 0 {
		topic = " [" + self.config.From + "]"
	}

	return fmt.Sprintf("%s%s%s%s", timenow, from, topic, content)
}

func (self Component) SendText(content string) error {
	content = self.generateContent(content)

	return gout.POST(self.config.Uri).SetJSON(gout.H{
		"msgtype": "text",
		"text": gout.H{
			"content":               content,
			"mentioned_list":        self.config.MentionedList,
			"mentioned_mobile_list": self.config.MentionedMobileList,
		},
	}).Debug(self.config.Debug).F().Retry().Attempt(self.config.Retry).WaitTime(time.Second).Do()
}

func (self Component) SendMarkDown(content string) error {
	content = self.generateContent(content)

	return gout.POST(self.config.Uri).SetJSON(gout.H{
		"msgtype": "markdown",
		"markdown": gout.H{
			"content":               content,
			"mentioned_list":        self.config.MentionedList,
			"mentioned_mobile_list": self.config.MentionedMobileList,
		},
	}).Debug(self.config.Debug).F().Retry().Attempt(self.config.Retry).WaitTime(time.Second).Do()
}

func logError(key string, err error) {
	log.Println(PackageName, ":error:"+key+":", err)
}
