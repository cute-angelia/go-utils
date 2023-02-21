package ibadger

import (
	"github.com/cute-angelia/go-utils/syntax/ijson"
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

func WithMemory(maxsize int64) Option {
	return func(c *Container) {
		c.config.Options.InMemory = true
		c.config.Options.MemTableSize = maxsize
		c.config.Options.Dir = ""
		c.config.Options.ValueDir = ""
	}
}

func WithPath(p string) Option {
	return func(c *Container) {
		c.config.Path = p
		// path
		c.config.Options.Dir = p
		c.config.Options.ValueDir = p
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config)
}
