package nuts

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
	return c
}

// Build
func (con *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(con)
	}
	return newComponent(con.config)
}

func WithDir(dir string) Option {
	return func(c *Container) {
		c.config.Options.Dir = dir
	}
}
