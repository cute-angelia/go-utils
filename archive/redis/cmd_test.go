package redis

import (
	"log"
	"testing"
)

func TestCmd(t *testing.T) {

	InitRedis("redis", "192.168.1.140", "6379", "")

	testKey := "STR:TEST:KEY:2"

	RedisCmder.Setex(testKey, 60, "1111")

	if v, err := RedisCmder.Get(testKey); err != nil {
		t.Error(err)
	} else {
		log.Println(v)
	}
}
