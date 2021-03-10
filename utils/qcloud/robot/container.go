package robot

import "fmt"

type Container struct {
	config *config
}

type Option func(c *Container)

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
	}
}

// Invoker ...
func Load(key string) *Container {
	c := DefaultContainer()
	c.config.Key = key
	c.config.Uri = fmt.Sprintf("%s?key=%s", c.config.Uri, c.config.Key)
	return c
}

// Build
func (con *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(con)
	}
	return newComponent(con.config)
}

func WithMentionedList(list []string) Option {
	return func(c *Container) {
		c.config.MentionedList = list
	}
}

func WithMentionedMobileList(list []string) Option {
	return func(c *Container) {
		c.config.MentionedMobileList = list
	}
}

func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}
