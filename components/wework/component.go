package wework

import (
	"fmt"
	"github.com/fastwego/wxwork/corporation"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
)

type Component struct {
	config *config
	app    *corporation.App
}

var Pools sync.Map

func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config
	return comp
}

func (c *Component) String() string {
	return ComponentName
}

// InitWeWork 初始化
func (c *Component) InitWeWork(redisClient *redis.Client) {
	log.Println(ComponentName, "InitWeWork", c.config.CorpId, c.config.AgentId)
	Corp := corporation.New(corporation.Config{
		Corpid: c.config.CorpId,
	})
	App := Corp.NewApp(corporation.AppConfig{
		AgentId: c.config.AgentId,
		Secret:  c.config.Secret,
	})
	App.SetAccessTokenCacheDriver(NewRedis(redisClient))

	_, loaded := Pools.LoadOrStore(c.config.AgentId, App)
	if loaded {
		log.Println("命中池")
	} else {
		log.Println("初始化")
	}
}

func GetApp(agentId string) (*corporation.App, error) {
	if v, ok := Pools.Load(agentId); ok {
		return v.(*corporation.App), nil
	} else {
		return nil, fmt.Errorf("%s: App不存在,请初始化配置", agentId)
	}
}
