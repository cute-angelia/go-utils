package iminio

import (
	"github.com/gotomicro/ego/core/econf"
	"log"
	"strings"
	"time"
)

type Option func(c *Container)

type Container struct {
	config *config
	name   string
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	// 两种方式，一种是 ego 的 config 加载，一种是option with 加载
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		panic("error parse config file")
		return c
	}

	c.name = key
	return c
}

func WithEndpoint(Endpoint string) Option {
	return func(c *Container) {
		c.config.Endpoint = Endpoint
	}
}
func WithAccesskeyId(AccesskeyId string) Option {
	return func(c *Container) {
		c.config.AccesskeyId = AccesskeyId
	}
}
func WithSecretaccessKey(SecretaccessKey string) Option {
	return func(c *Container) {
		c.config.SecretaccessKey = SecretaccessKey
	}
}
func WithDebug(Debug bool) Option {
	return func(c *Container) {
		c.config.Debug = Debug
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.Timeout = timeout
	}
}
func WithUseSSL(UseSSL bool) Option {
	return func(c *Container) {
		c.config.UseSSL = UseSSL
	}
}
func WithProxySocks5(ProxySocks5 string) Option {
	ProxySocks5 = strings.Replace(ProxySocks5, "socks5://", "", -1)

	return func(c *Container) {
		c.config.ProxySocks5 = ProxySocks5
	}
}

func WithReplaceMode(ReplaceMode int) Option {
	return func(c *Container) {
		c.config.ReplaceMode = ReplaceMode
	}
}

func WithReferer(Referer string) Option {
	return func(c *Container) {
		c.config.Referer = Referer
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if len(c.config.Endpoint) == 0 {
		log.Println("请初始化配置， 未能获取到配置信息")
	}

	// log.Println(PackageName, fmt.Sprintf("%+v", c.config))
	return newComponent(c.name, c.config)
}
