package idownload

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"strings"
	"time"
)

type Option func(c *Container)

type Container struct {
	config *config
	name   string
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	// 两种方式，一种是 ego 的 config 加载，一种是option with 加载
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

// Width  int // 图片-最小宽度
func WithWidth(Width int) Option {
	return func(c *Container) {
		c.config.Width = Width
	}
}

//Height int // 图片-最小高度
func WithHeight(Height int) Option {
	return func(c *Container) {
		c.config.Height = Height
	}
}

//Rename     bool   // 文件-是否重命名
func WithRename(Rename bool) Option {
	return func(c *Container) {
		c.config.Rename = Rename
	}
}

// Dest       string // 文件-下载路径
func WithDest(dest string) Option {
	return func(c *Container) {
		c.config.Dest = dest
	}
}

func WithDestAppend(dest string) Option {
	return func(c *Container) {
		c.config.Dest = c.config.Dest + "/" + dest
	}
}

// NamePrefix string // 文件-命名前缀
func WithNamePrefix(NamePrefix string) Option {
	return func(c *Container) {
		c.config.NamePrefix = NamePrefix
	}
}

//DefaultExt string // 文件-后缀属性，有些文件没有后缀
func WithDefaultExt(DefaultExt string) Option {
	return func(c *Container) {
		c.config.DefaultExt = DefaultExt
	}
}

//ProxyHttp   string // 代理 http://ip:port
func WithProxyHttp(ProxyHttp string) Option {
	return func(c *Container) {
		c.config.ProxyHttp = ProxyHttp
	}
}

//ProxySocks5 string // 代理 ip:port
func WithProxySocks5(ProxySocks5 string) Option {
	ProxySocks5 = strings.Replace(ProxySocks5, "socks5://", "", -1)

	return func(c *Container) {
		c.config.ProxySocks5 = ProxySocks5
	}
}

//Cookie      string // cookie
func WithCookie(Cookie string) Option {
	return func(c *Container) {
		c.config.Cookie = Cookie
	}
}

//UserAgent   string // user-agent
func WithUserAgent(UserAgent string) Option {
	return func(c *Container) {
		c.config.UserAgent = UserAgent
	}
}

// timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.Timeout = timeout
	}
}

// debug
func WithDebug(Debug bool) Option {
	return func(c *Container) {
		c.config.Debug = Debug
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.name, c.config, c.logger)
}
