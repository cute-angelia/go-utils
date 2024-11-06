package iredisV2

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
)

const name = "redis"

var redisMgrPools sync.Map

type RedisMgr struct {
	ctx    context.Context
	cancel context.CancelFunc
	opts   *options
	sfg    singleflight.Group // singleFlight

	alias string
}

func RedisInit(alias string, opts ...Option) *RedisMgr {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	if o.client == nil {
		// 如果指定了 MasterName 选项，那么返回哨兵模式的客户端 - FailoverClient
		// 如果选项 Addrs 的数量为两个或多个，那么返回集群模式的客户端 - ClusterClient
		// 否则就返回单节点的客户端
		o.client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        o.addrs,
			DB:           o.db,
			Username:     o.username,
			Password:     o.password,
			MaxRetries:   o.maxRetries,
			PoolSize:     o.poolSize,
			MinIdleConns: o.minIdleConns,
		})

		// 检测
		_, err := o.client.Ping(o.ctx).Result()
		if err != nil {
			panic(err)
		}
	}

	l := &RedisMgr{}
	l.ctx, l.cancel = context.WithCancel(o.ctx)
	l.opts = o

	// alias
	l.alias = alias

	_, loaded := redisMgrPools.LoadOrStore(alias, l)
	if loaded {
		l.printRedisPool("命中池", l.opts.client.PoolStats())
	} else {
		l.printRedisPool("初始化", l.opts.client.PoolStats())
	}

	return l
}

// GetRedisMgr 获取 redis 实例
func GetRedisMgr(alias string) (*RedisMgr, error) {
	if v, ok := redisMgrPools.Load(alias); ok {
		mgr := v.(*RedisMgr)
		return mgr, nil
	} else {
		return nil, fmt.Errorf("%s: reidsMgr不存在,请先初始化", alias)
	}
}

// GetRedis 获取 redis 实例
func GetRedis(alias string) (redis.UniversalClient, error) {
	if v, ok := redisMgrPools.Load(alias); ok {
		mgr := v.(*RedisMgr)
		return mgr.opts.client, nil
	} else {
		return nil, fmt.Errorf("%s: reids不存在,请先初始化", alias)
	}
}

func CloseAll() {
	redisMgrPools.Range(func(key, value interface{}) bool {
		mgr := value.(*RedisMgr)
		mgr.opts.client.Close()
		redisMgrPools.Delete(key)
		return true
	})
}

func (c *RedisMgr) printRedisPool(msg string, stats *redis.PoolStats) {
	log.Printf("[%s] Alias:%s, Sserver=%s Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d %s \n",
		name,
		c.opts.addrs,
		c.alias,
		stats.Hits,
		stats.Misses,
		stats.Timeouts,
		stats.TotalConns,
		stats.IdleConns,
		stats.StaleConns,
		msg,
	)
}
