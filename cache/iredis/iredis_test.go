package iredis

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	InitRedis("cache", "47.99.166:6379", "")
	rdb, _ := GetRedis("cache")
	ctx := context.Background()

	rdb.Set(ctx, "abc", "world", time.Second*10)

	val, err := rdb.Get(ctx, "abc").Result()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("key", val)
	}
}
