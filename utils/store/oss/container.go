package oss

import (
	"fmt"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"log"
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

func WithAccessKeyId(AccessKeyId string) Option {
	return func(c *Container) {
		c.config.AccessKeyId = AccessKeyId
	}
}
func WithAccessKeySecret(AccessKeySecret string) Option {
	return func(c *Container) {
		c.config.AccessKeySecret = AccessKeySecret
	}
}
func WithEndpoint(Endpoint string) Option {
	return func(c *Container) {
		c.config.Endpoint = Endpoint
	}
}
func WithBucketName(BucketName string) Option {
	return func(c *Container) {
		c.config.BucketName = BucketName
	}
}
func WithBucketHost(BucketHost string) Option {
	return func(c *Container) {
		c.config.BucketHost = BucketHost
	}
}
func WithDebug(Debug bool) Option {
	return func(c *Container) {
		c.config.Debug = Debug
	}
}

// Build ...
func (c *Container) MustBuild(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if len(c.config.BucketName) == 0 ||
		len(c.config.AccessKeyId) == 0 ||
		len(c.config.AccessKeySecret) == 0 {
		log.Println(PackageName, fmt.Sprintf("ERROR: %+v", c.config))
	}
	return newComponent(c.name, c.config, c.logger)
}
