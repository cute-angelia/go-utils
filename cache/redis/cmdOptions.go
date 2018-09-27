package redis

import "context"

type CmdOptions struct {
	Name string
	Ctx  context.Context
}
type CmdOption func(*CmdOptions)

func newCmdOption(opts ...CmdOption) CmdOptions {
	var sopt CmdOptions
	for _, opt := range opts {
		opt(&sopt)
	}
	if sopt.Name == "" {
		sopt.Name = "redis"
	}
	return sopt
}

func WithName(name string) CmdOption {
	return func(options *CmdOptions) {
		options.Name = name
	}
}

func WithContext(ctx context.Context) CmdOption {
	return func(options *CmdOptions) {
		options.Ctx = ctx
	}
}
