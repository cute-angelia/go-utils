package pool

import (
	"testing"
	"time"
)

func TestAnts(t *testing.T) {
	p := NewPoolAnts(100)
	p.InitTask(10)
	for i := 0; i < 10; i++ {

		// pool 一些都会重新赋值
		z := i
		p.SubmitTask(func() {
			t.Log("task:", z)
		})
	}
	t.Log("RunningTask", p.RunningTask())
	defer p.stop()

	time.Sleep(time.Second*1)
}
