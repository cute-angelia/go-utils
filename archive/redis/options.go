package redis

import (
	"time"
)

const (
	defaultReadyTimeout = time.Second * 3
)

type Option func(options *Options)

type Options struct {
	Name           string
	DialTimeout    time.Duration
	MaxDialRetries int
	MaxPerHost     int
	IdleTime       time.Duration
	ReadyTimeout   time.Duration

	Addr string
	Auth string
}

func newOption(opts ...Option) Options {
	var sopt Options
	for _, opt := range opts {
		opt(&sopt)
	}

	if sopt.DialTimeout == 0 {
		sopt.DialTimeout = time.Millisecond * 150
	}

	if sopt.MaxDialRetries == 0 {
		sopt.MaxDialRetries = 3
	}

	if sopt.MaxPerHost == 0 {
		sopt.MaxPerHost = 20
	}

	if sopt.IdleTime == 0 {
		sopt.IdleTime = time.Minute * 7
	}

	if sopt.ReadyTimeout == 0 {
		sopt.ReadyTimeout = defaultReadyTimeout
	}

	return sopt
}

func DialTimeout(d time.Duration) Option {
	return func(options *Options) {
		options.DialTimeout = d
	}
}

func MaxDialRetries(i int) Option {
	return func(options *Options) {
		options.MaxDialRetries = i
	}
}

func Name(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func MaxPerHost(n int) Option {
	return func(options *Options) {
		options.MaxPerHost = n
	}
}

func IdleTime(d time.Duration) Option {
	return func(options *Options) {
		options.IdleTime = d
	}
}

func MakeConfig(addr string, auth string) Option {
	return func(options *Options) {
		options.Addr = addr
		options.Auth = auth
	}
}
