package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"strings"
)

type Pool interface {
	Get(name string) redigo.Conn
	GetWithCtx(ctx context.Context, name string) redigo.Conn
}

func New(name string, opts ...Option) Pool {
	options := newOption(opts...)

	p := &pool{
		opts: options,
		name: name,
	}

	return p
}

type pool struct {
	opts Options
	name string

	mu    sync.RWMutex
	pools sync.Map // map[string]*redigo.Pool
}

// 获取 Get
func (r *pool) Get(name string) redigo.Conn {
	conn, err := r.getConn(name)
	if err != nil {
		return &errorConnection{err: err}
	}
	return conn
}

// 获取 GetWithCtx
func (r *pool) GetWithCtx(ctx context.Context, name string) redigo.Conn {
	conn, err := r.getConn(name)
	if err != nil {
		return &errorConnection{err: err, ctx: ctx}
	}
	conn.ctx = ctx
	return conn
}

// 获取 conn
func (r *pool) getConn(name string) (*conn, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, exist := r.pools.Load(name); !exist {
		p2 := newPool(&r.opts)
		r.pools.Store(name, p2)
		return &conn{Conn: p2.Get()}, nil
	} else {
		return &conn{Conn: p.(*redigo.Pool).Get()}, nil
	}
}

// new pool
func newPool(opt *Options) *redigo.Pool {
	dialFunc := func() (redigo.Conn, error) {
		for i := opt.MaxDialRetries; i > 0; i-- {
			if strings.Index(opt.Addr, "://") == -1 {
				if c, err := redigo.Dial("tcp", opt.Addr, redigo.DialConnectTimeout(opt.DialTimeout)); err == nil {
					if len(opt.Auth) > 0 {
						if _, err := c.Do("AUTH", opt.Auth); err != nil {
							c.Close()
							return nil, err
						}
					}
					return c, nil
				}
			} else {
				if c, err := redigo.DialURL(opt.Addr, redigo.DialConnectTimeout(opt.DialTimeout)); err == nil {
					if len(opt.Auth) > 0 {
						if _, err := c.Do("AUTH", opt.Auth); err != nil {
							c.Close()
							return nil, err
						}
					}
					return c, nil
				}
			}

		}
		return nil, fmt.Errorf("dial %s failed", opt.Addr)
	}

	return &redigo.Pool{
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		Dial:        dialFunc,
		MaxIdle:     opt.MaxPerHost,
		IdleTimeout: time.Minute * 7,
	}
}
