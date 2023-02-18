package iredis

import (
	"github.com/cute-angelia/go-utils/v2/components/ijson"
	"github.com/spf13/viper"
	"log"
)

type Option func(c *Container)

type Container struct {
	config *config
}

// Load viper 加载 配置
func Load(key string) *Component {
	iconfig := DefaultConfig()
	configData := viper.GetStringMap(key)
	jsonstr, _ := ijson.Marshal(configData)
	if err := ijson.Unmarshal(jsonstr, &iconfig); err != nil {
		log.Println(err)
	}
	// log.Println(ijson.Pretty(iconfig))
	return newComponent(iconfig)
}

// New options 模式
func New(options ...Option) *Component {
	c := &Container{
		config: DefaultConfig(),
	}
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config)
}

func WithName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
	}
}
func WithServer(server string) Option {
	return func(c *Container) {
		c.config.Server = server
	}
}
func WithPassword(password string) Option {
	return func(c *Container) {
		c.config.Password = password
	}
}
func WithDbIndex(dbIndex int) Option {
	return func(c *Container) {
		c.config.DbIndex = dbIndex
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config)
}
