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

// 集群
func (c Component) initRedisUniversal() {
	// redis.NewUniversalClient()
}

func (c Component) initRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.config.Server,
		Password: c.config.Password, // no password set
		DB:       c.config.DbIndex,  // use default DB

		//连接池容量及闲置连接数量
		//PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		//MinIdleConns: 10, // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；

		//超时
		//DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		//ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		//WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		//PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		//IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		//IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		//MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		//MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		//MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		//MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

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
