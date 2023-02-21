package redis

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestDo(t *testing.T) {

	p := New("redis", MakeConfig("192.168.1.140:6379", ""))
	conn := p.Get("redis")

	for i := 0; i < 1; i++ {

		if _, err := conn.Do("SET", "STR:TEST:KEY1", 1); err != nil {
			t.Error(err)
		}

		if i, err := redis.Int(conn.Do("GET", "STR:TEST:KEY1")); err != nil {
			t.Error(err)
		} else if i != 1 {
			t.Error("val incorrect")
		}

		if _, err := conn.Do("DEL", "STR:TEST:KEY1", 1); err != nil {
			t.Error(err)
		}
		conn.Close()
	}
}

func BenchmarkGet(b *testing.B) {
	p := New("redis", MakeConfig("192.168.1.140:6379", ""))
	for i := 0; i < b.N; i++ {
		conn := p.Get("redis")
		conn.Close()
	}
}
