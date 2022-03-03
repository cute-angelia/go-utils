package ffmpeg

import (
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

func Load() *Container {
	c := DefaultContainer()
	return c
}

func WithFfmpegPath(ffmpegPath string) Option {
	return func(c *Container) {
		c.config.FfmpegPath = ffmpegPath
	}
}

func WithTimeOut(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.Timeout = timeout
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.name, c.config, c.logger)
}
