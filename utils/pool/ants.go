package pool

import "github.com/panjf2000/ants"

var AntsPool *ants.Pool

func NewPoolAnts(num int) *ants.Pool {
	if p, err := ants.NewPool(num, ants.WithPreAlloc(true)); err != nil {
		panic(err)
	} else {
		AntsPool = p
		return p
	}
}
