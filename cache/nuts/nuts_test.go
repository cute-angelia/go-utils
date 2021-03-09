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
