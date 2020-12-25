package iredis2

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var RedisPools *pool

type pool struct {
	Pool map[string]*redis.Client
}

func NewRedisPool() *pool {
	return &pool{}
}

func (self *pool) GetRedis(name string) (*redis.Client, error) {
	if _, ok := self.Pool[name]; ok {
		return self.Pool[name], nil
	} else {
		return nil, fmt.Errorf("%s: reids不存在初始化配置", name)
	}
}

func GetRdb(name string) *redis.Client {
	rdb, _ := RedisPools.GetRedis(name)
	return rdb
}

func Init(name string, server string, password string) *pool {
	if RedisPools == nil {
		RedisPools = &pool{}
	}

	if RedisPools.Pool == nil {
		RedisPools.Pool = make(map[string]*redis.Client)
	}

	if _, ok := RedisPools.Pool[name]; ok {

		// 打印
		printRedisPool(name, RedisPools.Pool[name].PoolStats())

		// 赋值
		return RedisPools
	} else {
		RedisPools.Pool[name] = redis.NewClient(&redis.Options{
			Addr:     server,
			Password: password, // no password set
			DB:       0,        // use default DB
		})

		// 打印
		printRedisPool(name, RedisPools.Pool[name].PoolStats())
		return RedisPools
	}
}

func Close() {
	for _, i2 := range RedisPools.Pool {
		i2.Close()
	}
}

func printRedisPool(name string, stats *redis.PoolStats) {
	log.Printf("Redis初始化：Name:%s, Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n", name,
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}
