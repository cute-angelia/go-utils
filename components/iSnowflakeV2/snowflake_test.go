package iSnowflakeV2

import (
	"log"
	"testing"
	"time"
)

func TestGen(t *testing.T) {
	isnow := NewSnowflake(WithWorkerId(1))
	m := make(map[int]int64)
	for i := 0; i < 10000; i++ {
		id := isnow.GetId()
		if _, ok := m[i]; ok {
			log.Panicln("error: 重复了")
		} else {
			m[i] = id
		}
	}
	log.Println(m, time.Now().Unix())

}
