package pool

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log"
	"sync"
)

type AntsPool struct {
	pool *ants.Pool
	wg   *sync.WaitGroup
}

// num 协程数量 , 预先分配内存
func MustNewPoolAnts(num int, withPreAlloc bool) *AntsPool {
	if num <= 0 {
		panic(fmt.Errorf("需要制定协程数量"))
	}
	ap := AntsPool{}
	var wg sync.WaitGroup
	// 预先分配内存 WithPreAlloc
	if p, err := ants.NewPool(num, ants.WithPreAlloc(withPreAlloc)); err != nil {
		panic(err)
	} else {
		ap.pool = p
		ap.wg = &wg
	}
	return &ap
}

//// 步骤 1  任务数量
//func (t *AntsPool) WgAdd() *AntsPool {
//	t.wg.Add(1)
//	return t
//}

// 步骤 2
func (t *AntsPool) SubmitTask(fc func()) error {
	t.wg.Add(1)
	if err := t.pool.Submit(fc); err != nil {
		log.Println(err)
		t.wg.Done()
		return err
	} else {
		t.wg.Done()
		return nil
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
