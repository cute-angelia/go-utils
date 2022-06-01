package iredis

import (
	"context"
	"github.com/cute-angelia/go-utils/utils/conf"
	"log"
	"testing"
	"time"
)

func startConnectRedis() {
	conf.LoadConfigFile("./config.toml")
	Load("redis")
}

func TestRedis(t *testing.T) {

	// start
	startConnectRedis()

	redisKey := "abc"

	rdb, _ := GetRedis("cache")
	ctx := context.Background()
	rdb.Set(ctx, redisKey, "world", time.Minute*1)
	val, err := rdb.Get(ctx, redisKey).Result()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(redisKey, val)
	}

	t.Log("Exists:", rdb.Exists(ctx, redisKey+"x").Val())
}
