package idownload

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"runtime"
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

// WithProxyHttp   string // 代理 http://ip:port
func WithProxyHttp(ProxyHttp string) Option {
	return func(c *Container) {
		c.config.ProxyHttp = ProxyHttp
	}
}

// WithProxySocks5 string // 代理 ip:port
func WithProxySocks5(ProxySocks5 string) Option {
	ProxySocks5 = strings.Replace(ProxySocks5, "socks5://", "", -1)

	return func(c *Container) {
		c.config.ProxySocks5 = ProxySocks5
	}
}

// WithCookie  string // cookie
func WithCookie(Cookie string) Option {
	return func(c *Container) {
		c.config.Cookie = Cookie
	}
}

// WithUserAgent user-agent
func WithUserAgent(UserAgent string) Option {
	return func(c *Container) {
		c.config.UserAgent = UserAgent
	}
}

func WithUserAgentRandon(random bool) Option {
	return func(c *Container) {
		c.config.UseRandomUserAgent = random
	}
}

func WithUserAgentRandonMobile(random bool) Option {
	return func(c *Container) {
		c.config.UseRandomUserAgentMobile = random
	}
}

func WithReferer(Referer string) Option {
	return func(c *Container) {
		c.config.Referer = Referer
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.Timeout = timeout
	}
}

func WithDebug(Debug bool) Option {
	return func(c *Container) {
		c.config.Debug = Debug
	}
}

func WithResume(resume bool) Option {
	return func(c *Container) {
		c.config.Resume = resume
	}
}

func WithAuthorization(authorization string) Option {
	return func(c *Container) {
		c.config.Authorization = authorization
	}
}

func WithConcurrency(concurrency int) Option {
	return func(c *Container) {
		c.config.Concurrency = concurrency
	}
}
func WithConcurrencyCpu() Option {
	concurrency := runtime.NumCPU()
	if concurrency >= 2 {
		concurrency = concurrency / 2
	}
	return func(c *Container) {
		c.config.Concurrency = concurrency
	}
}
func WithRetryAttempt(retryAttempt int) Option {
	return func(c *Container) {
		c.config.RetryAttempt = retryAttempt
	}
}
func WithRetryWaitTime(retryWaitTime time.Duration) Option {
	return func(c *Container) {
		c.config.RetryWaitTime = retryWaitTime
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if c.config.UseRandomUserAgent {
		c.config.UserAgent = RandomUserAgent()
	}
	if c.config.UseRandomUserAgentMobile {
		c.config.UserAgent = RandomMobileUserAgent()
	}

	return newComponent(c.name, c.config, c.logger)
}
