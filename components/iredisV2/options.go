package iredisV2

import (
	"context"
	"github.com/go-redis/redis/v8"
	"runtime"
)

const (
	defaultAddr         = "127.0.0.1:6379"
	defaultDB           = 0
	defaultMaxRetries   = 3
	defaultPoolSize     = 10 // 默认连接池大小等于 runtime.GOMAXPROCS(cpu个数) * 10
	defaultMinIdleConns = 1  // 启动时，创建最小空闲连接数
)

type Option func(o *options)

type options struct {
	ctx context.Context

	// 客户端连接地址
	// 内建客户端配置，默认为[]string{"127.0.0.1:6379"}
	addrs []string

	// 数据库号
	// 内建客户端配置，默认为0 范围:0~15
	db int

	// 用户名
	// 内建客户端配置，默认为空
	username string

	// 密码
	// 内建客户端配置，默认为空
	password string

	// 最大重试次数
	// 内建客户端配置，默认为3次
	maxRetries int

	// 客户端
	// 外部客户端配置，存在外部客户端时，优先使用外部客户端，默认为nil
	client redis.UniversalClient

	poolSize int

	minIdleConns int
}

func defaultOptions() *options {
	return &options{
		ctx:          context.Background(),
		addrs:        []string{defaultAddr},
		db:           defaultDB,
		maxRetries:   defaultMaxRetries,
		poolSize:     defaultPoolSize * runtime.GOMAXPROCS(runtime.NumCPU()),
		minIdleConns: defaultMinIdleConns,
		username:     "",
		password:     "",
	}
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithAddrs 设置连接地址
func WithAddrs(addrs ...string) Option {
	return func(o *options) { o.addrs = addrs }
}

// WithDB 设置数据库号
func WithDB(db int) Option {
	return func(o *options) { o.db = db }
}

// WithUsername 设置用户名
func WithUsername(username string) Option {
	return func(o *options) { o.username = username }
}

// WithPassword 设置密码
func WithPassword(password string) Option {
	return func(o *options) { o.password = password }
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) Option {
	return func(o *options) { o.maxRetries = maxRetries }
}

// WithClient 设置外部客户端
func WithClient(client redis.UniversalClient) Option {
	return func(o *options) { o.client = client }
}

func WithPoolSize(poolSize int) Option {
	return func(o *options) { o.poolSize = poolSize }
}

func WithMinIdleConns(minIdleConns int) Option {
	return func(o *options) { o.minIdleConns = minIdleConns }
}
