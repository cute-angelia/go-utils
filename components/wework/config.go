package wework

const ComponentName = "Component.WeWork"

// 企业微信SDK
type config struct {
	CorpId  string
	AgentId string
	Secret  string
}

func DefaultConfig() *config {
	c := config{}
	return &c
}
