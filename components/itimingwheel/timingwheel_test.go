package itimingwheel

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestDemo(t *testing.T) {

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
			}
		}
	}()

	// do
	var count uint64
	tw, err := NewTimingWheel(time.Second, 20, func(key, value interface{}) {
		log.Println(key, value, "key-value")
		job(&count)
	})
	if err != nil {
		log.Fatal(err)
	}

	defer tw.Stop()
	for i := 0; i < 10; i++ {
		tw.SetTimer(i, i, time.Second*1)
	}

	tw.SetTimer("ok", "value", time.Second*10)
	tw.MoveTimer("ok", time.Second*4) // 修改延迟到 4 秒，就执行

	tw.run()
}

func job(count *uint64) {
	v := atomic.AddUint64(count, 1)
	if v%10000 == 0 {
		fmt.Println(v)
	}
}
