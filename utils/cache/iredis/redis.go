package iredis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
)

var RedisPools sync.Map

// 初始化 Redis
func InitRedis(name string, server string, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     server,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	_, loaded := RedisPools.LoadOrStore(name, redisClient)

	if loaded {
		printRedisPool(name+"——命中池", redisClient.PoolStats())
	} else {
		printRedisPool(name+"——初始化", redisClient.PoolStats())
	}
	return redisClient
}

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

func printRedisPool(name string, stats *redis.PoolStats) {
	log.Printf("Redis初始化：Name:%s, Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n", name,
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}
