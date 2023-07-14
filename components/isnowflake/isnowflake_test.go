package isnowflake

import (
	"log"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	isnlow := New(WithWorkerId(1))
	for i := 0; i < 100; i++ {
		log.Println(isnlow.GetId())
	}
}
