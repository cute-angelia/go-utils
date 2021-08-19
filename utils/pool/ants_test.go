package pool

import (
	"log"
	"testing"
	"time"
)

func TestAnts(t *testing.T) {

	p := NewPoolAnts(100)
	defer p.Stop()

	generate(p)

	time.Sleep(time.Second * 1)
}

func generate(p *AntsPool) {

	as := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	p.InitTask(10 + len(as))

	for _, i2 := range as {
		// notice 必须重新赋值
		newi2 := i2
		p.SubmitTask(func() {
			dosomethings("--> stp1", newi2)
		})
	}

	for i := 0; i < 10; i++ {
		// pool 一些都会重新赋值
		z := i
		p.SubmitTask(func() {
			dosomethings("stp2", z)
		})
	}
	log.Println("RunningTask", p.RunningTask())
}

func dosomethings(flag string, z int) {
	log.Println("task:", flag, z)
}
