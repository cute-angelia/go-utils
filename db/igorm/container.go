package igorm

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Option func(c *Container)

type Container struct {
	config *config
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
	}
}

func Load(dbName string) *Container {
	c := DefaultContainer()
	c.config.DbName = dbName
	return c
}

// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量(10)
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *Container) {
		c.config.MaxIdleConns = maxIdleConns
	}
}

// SetMaxOpenConns 设置打开数据库连接的最大数量(100)
func WithMaxOpenConnss(maxOpenConns int) Option {
	return func(c *Container) {
		c.config.MaxOpenConns = maxOpenConns
	}
}

func WithMaxLifetime(maxLifetime time.Duration) Option {
	return func(c *Container) {
		c.config.MaxLifetime = maxLifetime
	}
}

func WithLogDebug(logDebug bool) Option {
	return func(c *Container) {
		c.config.LogDebug = logDebug
	}
}

func WithLogger(logger gorm.Logger) Option {
	return func(c *Container) {
		c.config.Logger = logger
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	return newComponent(c.config)
}
