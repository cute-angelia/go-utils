package apicache

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
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

func WithTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.Timeout = timeout
	}
}
func WithOnlyToday(onlyTodya bool) Option {
	return func(c *Container) {
		c.config.OnlyToday = onlyTodya
	}
}
func WithCacheKey(cacheKey string) Option {
	return func(c *Container) {
		c.config.CacheKey = cacheKey
	}
}
func WithDbName(dbName string) Option {
	return func(c *Container) {
		c.config.DbName = dbName
	}
}
func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}
func WithDeleteKey(deleteKey string) Option {
	return func(c *Container) {
		c.config.DeleteKey = deleteKey
	}
}

// Build ...
func (c *Container) MustBuild(dbName string, cacheKey string, options ...Option) *Component {

	// Must
	c.config.DbName = dbName
	c.config.CacheKey = cacheKey

	for _, option := range options {
		option(c)
	}
	// log.Println(PackageName, fmt.Sprintf("%+v", c.config))

	if len(c.config.CacheKey) == 0 {
		c.logger.Error("cachekey is empty")
		panic("cachekey is empty")
	}
	if len(c.config.DbName) == 0 {
		c.logger.Error("DbName is empty")
		panic("DbName is empty")
	}

	if c.config.Debug {
		c.logger.SetLevel(elog.DebugLevel)
	}

	return newComponent(c.name, c.config, c.logger)
}
