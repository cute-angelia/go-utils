package iredisV2

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func (c *RedisMgr) Expire(key string, timeDur time.Duration) {
	c.opts.client.Expire(c.ctx, key, timeDur)
}

func (c *RedisMgr) Set(key string, val string, expire time.Duration) error {
	return c.opts.client.SetEX(c.ctx, key, val, expire).Err()
}

// Get
func (c *RedisMgr) Get(key string) string {
	val, err, _ := c.sfg.Do(key, func() (interface{}, error) {
		val, err := c.opts.client.Get(c.ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println("RedisMgr Get Error1", err)
			return "", err
		}
		return val, nil
	})
	if err != nil {
		log.Println("RedisMgr Get Error2", err)
		return ""
	}
	return val.(string)
}

func (c *RedisMgr) HSet(key string, field string, value string) error {
	return c.opts.client.HSet(c.ctx, key, field, value).Err()
}

func (c *RedisMgr) HGet(key string, field string) string {
	val, err, _ := c.sfg.Do(key+field, func() (interface{}, error) {
		val, err := c.opts.client.HGet(c.ctx, key, field).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println("RedisMgr Get Error1", err)
			return "", err
		}
		return val, nil
	})
	if err != nil {
		log.Println("RedisMgr Get Error2", err)
		return ""
	}
	return val.(string)
}

func (c *RedisMgr) HMSet(key string, data map[string]interface{}) error {
	return c.opts.client.HMSet(c.ctx, key, data).Err()
}

func (c *RedisMgr) HMGet(key string, fields ...string) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(fields))
	if vals, err := c.opts.client.HMGet(c.ctx, key, fields...).Result(); err != nil {
		return result, err
	} else {
		for i, val := range fields {
			result[val] = vals[i]
		}
		return result, nil
	}
}

func (c *RedisMgr) HExists(key string, field string) (bool, error) {
	return c.opts.client.HExists(c.ctx, key, field).Result()
}

func (c *RedisMgr) HDel(key string, field string) (bool, error) {
	_, err := c.opts.client.HDel(c.ctx, key, field).Result()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// LRange 分页获取数据
func (c *RedisMgr) LRange(key string, page int64, perpage int64) ([]string, error) {
	start := (page - 1) * perpage
	stop := start + perpage
	val, err, _ := c.sfg.Do(fmt.Sprintf("%s_%d_%d", key, page, perpage), func() (interface{}, error) {
		val, err := c.opts.client.LRange(c.ctx, key, start, stop).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return []string{}, err
		}
		return val, nil
	})
	if err != nil {
		log.Println("RedisMgr Get Error2", err)
		return []string{}, err
	}
	return val.([]string), nil
}

// LTrimLimit 只保留N个数据
func (c *RedisMgr) LTrimLimit(key string, n int64) (string, error) {
	if n >= 1 {
		n--
	}
	return c.opts.client.LTrim(c.ctx, key, 0, n).Result()
}
