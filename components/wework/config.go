package wework

import "github.com/go-redis/redis/v8"

const ComponentName = "Component.WeWork"

// 企业微信SDK
type config struct {
	CorpId      string
	AgentId     string
	Secret      string
	redisClient *redis.Client
}

func DefaultConfig() *config {
	c := config{}
	return &c
}
