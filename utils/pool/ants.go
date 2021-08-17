package pool

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

type AntsPool struct {
	pool    *ants.Pool
	wg      *sync.WaitGroup
	taskNum int
}

func NewPoolAnts(num int) *AntsPool {
	if num <= 0 {
		panic(fmt.Errorf("需要制定协程数量"))
	}

	ap := AntsPool{}
	// 预先分配内存 WithPreAlloc
	if p, err := ants.NewPool(num, ants.WithPreAlloc(true)); err != nil {
		panic(err)
	} else {
		ap.pool = p
	}

	return &ap
}

// 步骤 1
func (t *AntsPool) InitTask(taskNum int) {
	var wg sync.WaitGroup

	// 检查数量
	t.taskNum = taskNum
	if taskNum > 0 {
		wg.Add(taskNum)
	}

	t.wg = &wg
}

// 步骤 2
func (t *AntsPool) SubmitTask(fc func()) error {
	if t.taskNum <= 0 {
		return fmt.Errorf("you should InitTask, 需要指定任务数量")
	} else {
		err := t.pool.Submit(fc)
		t.wg.Done()
		return err
	}
}

// 步骤 3
func (t *AntsPool) RunningTask() int {
	t.wg.Wait()
	return t.pool.Running()
}

// 步骤 4
func (t *AntsPool) Stop() {
	t.pool.Release()
}
