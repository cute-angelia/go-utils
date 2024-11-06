package iredisV2

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {

	// 初始化
	RedisInit("cache", WithAddrs("127.0.0.1:6379"))

	rdb, _ := GetRedis("cache")
	rdbMgr, _ := GetRedisMgr("cache")
	ctx := context.Background()

	log.Println("--------------- 测试 set get ---------------")
	redisKey := "abc"
	rdb.Set(ctx, redisKey, "world", time.Minute*1)
	val, err := rdb.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println("错误", err)
	} else {
		log.Println(redisKey, val)
	}
	t.Log("Exists:", rdb.Exists(ctx, redisKey).Val())

	log.Println("--------------- 测试 hmset hmget ---------------")
	hmsetkey := "hmsetkey_tset"
	rdbMgr.HMSet(hmsetkey, map[string]interface{}{
		"user1": "user1-ok",
		"user2": "user2-ok",
		"user4": "user4-ok",
	})
	rdbMgr.Expire(hmsetkey, time.Minute)
	if result, err := rdbMgr.HMGet(hmsetkey, "user1", "user2", "user3", "user4"); err != nil {
		log.Println(err)
	} else {
		log.Println(result)
		for s, i := range result {
			if i != nil {
				log.Println(s, i.(string))
			}
		}
	}

	log.Println("--------------- 测试 hset hget ---------------")
	hsetkey := "hsetkey_tset"
	log.Println(rdbMgr.HGet(hsetkey, "name"))
	log.Println("hset err:", rdbMgr.HSet(hsetkey, "name", "abc"))
	log.Println(rdbMgr.HGet(hsetkey, "name"))

	log.Println("--------------- 测试 LTRIM ---------------")
	htrimkey := "hltrim_tset"
	rdb.Del(ctx, htrimkey)
	rdb.RPush(ctx, htrimkey, "user1", "user2", "user3", "user4")
	rdb.LTrim(ctx, htrimkey, 0, 100) // 只保留指定区间内的元素，不在指定区间之内的元素都将被删除
	log.Println(rdbMgr.LTrimLimit(htrimkey, 3))
	log.Println(rdb.LRange(ctx, htrimkey, 0, -1).Result())
}
