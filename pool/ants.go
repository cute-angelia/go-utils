package pool

import "github.com/panjf2000/ants"

var AntsPool *ants.Pool

func NewPoolAnts(num int) *ants.Pool {
	p, _ := ants.NewPool(num, ants.WithPreAlloc(true))
	return p
}