package iredis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
)

const PackageName = "component.iredis"

var RedisPools sync.Map

type Component struct {
	config *config
}

// newComponent ...
func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config
	// 初始化 redis
	comp.initRedis()
	return comp
}

func (c Component) initRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.config.Server,
		Password: c.config.Password, // no password set
		DB:       c.config.DbIndex,  // use default DB
	})
	_, loaded := RedisPools.LoadOrStore(c.config.Name, redisClient)

	if loaded {
		c.printRedisPool("命中池", redisClient.PoolStats())
	} else {
		c.printRedisPool("初始化", redisClient.PoolStats())
	}
	return redisClient
}

func (c Component) printRedisPool(msg string, stats *redis.PoolStats) {
	log.Printf("[%s] Name:%s, Sserver=%s Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d %s \n",
		PackageName,
		c.config.Name,
		c.config.Server,
		stats.Hits,
		stats.Misses,
		stats.Timeouts,
		stats.TotalConns,
		stats.IdleConns,
		stats.StaleConns,
		msg,
	)
}

// GetRedisClint 获取 redis 实例
func GetRedisClient(name string) *redis.Client {
	if v, ok := RedisPools.Load(name); ok {
		return v.(*redis.Client)
	} else {
		return nil
	}
}

// GetRedis 获取 redis 实例
func GetRedis(name string) (*redis.Client, error) {
	if v, ok := RedisPools.Load(name); ok {
		return v.(*redis.Client), nil
	} else {
		return nil, fmt.Errorf("%s: reids不存在,请初始化配置", name)
	}
}

func Close() {
	RedisPools.Range(func(key, value interface{}) bool {
		value.(*redis.Client).Close()
		return true
	})
}
