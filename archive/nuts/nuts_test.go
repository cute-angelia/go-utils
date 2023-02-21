package nuts

import (
	"log"
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	nuts := Load("test").Build(WithDir("/tmp/nutsdb"))
	bucket := "test"

	defer func() {
		nuts.Merge()
		nuts.Close()
	}()

	nuts.Set(bucket, "hello", "world", 3)

	for {
		time.Sleep(time.Second)
		log.Println("hello:", nuts.Get(bucket, "hello"))
		log.Println("hello1:", nuts.Get(bucket, "hello1"))
	}
}

func TestLocker(t *testing.T) {
	nuts := Load("test").Build(WithDir("/tmp/nutsdb"))
	bucket := "test"

	for {
		time.Sleep(time.Second)
		opts := NewLockerOpt(WithLimit(10), WithToday(true))
		if nuts.IsNotLockedInLimit(bucket, "hello", 86400, opts) {
			log.Println("i in")
		}
	}
}

func TestIncr(t *testing.T) {
	nuts := Load("test").Build(WithDir("/tmp/nutsdb"))
	bucket := "test"

	for {
		nuts.Incr(bucket, "test", "1", 100)

		log.Println(nuts.Get(bucket, "test"))
	}
}
