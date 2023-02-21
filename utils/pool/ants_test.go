package pool

import (
	"log"
	"testing"
	"time"
)

func TestAnts(t *testing.T) {
	p := MustNewPoolAnts(2, false)
	defer p.Stop()

	as := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i2 := range as {
		// notice 必须重新赋值
		newi2 := i2
		p.SubmitTask(func() {
			log.Println("--> ", newi2)
		})
	}
	log.Println("RunningTask", p.RunningTask())
	<-time.After(time.Second * 1)
}
