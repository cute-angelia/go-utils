package loggerV3

import "github.com/rs/zerolog"

type Option func(c *Container)

type Container struct {
	config *config
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
	}
}

func Load() *Container {
	return DefaultContainer()
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config)
}

func WithProject(project string) Option {
	return func(c *Container) {
		c.config.project = project
	}
}

func WithIsOnline(isOnline bool) Option {
	return func(c *Container) {
		c.config.isOnline = isOnline
	}
}

func WithMaxSize(maxSize int) Option {
	return func(c *Container) {
		c.config.maxSize = maxSize
	}
}
func WithMaxBackups(maxBackups int) Option {
	return func(c *Container) {
		c.config.maxBackups = maxBackups
	}
}
func WithMaxAge(maxAge int) Option {
	return func(c *Container) {
		c.config.maxAge = maxAge
	}
}
func WithEveryday(everyday bool) Option {
	return func(c *Container) {
		c.config.everyday = everyday
	}
}

func WithLevel(level zerolog.Level) Option {
	return func(c *Container) {
		c.config.level = level
	}
}

func WithFileJson(fileJson bool) Option {
	return func(c *Container) {
		c.config.fileJson = fileJson
	}
}
