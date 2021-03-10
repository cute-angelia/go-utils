package robot

import (
	"github.com/guonaihong/gout"
	"log"
)

type Component struct {
	config *config
}

func newComponent(cfg *config) *Component {
	return &Component{
		config: cfg,
	}
}

func (self Component) SendText(content string) {
	gout.POST(self.config.Uri).SetJSON(gout.H{
		"msgtype": "text",
		"text": gout.H{
			"content":               content,
			"mentioned_list":        self.config.MentionedList,
			"mentioned_mobile_list": self.config.MentionedMobileList,
		},
	}).Debug(self.config.Debug).Do()
}

func (self Component) SendMarkDown(content string) {
	gout.POST(self.config.Uri).SetJSON(gout.H{
		"msgtype": "markdown",
		"markdown": gout.H{
			"content":               content,
			"mentioned_list":        self.config.MentionedList,
			"mentioned_mobile_list": self.config.MentionedMobileList,
		},
	}).Debug(self.config.Debug).Do()
}

func logError(key string, err error) {
	log.Println(PackageName, ":error:"+key+":", err)
}
