package iupload

import (
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

func WithReplaceMode(ReplaceMode int) Option {
	return func(c *Container) {
		c.config.ReplaceMode = ReplaceMode
	}
}

func (c *Container) New(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config)
}
