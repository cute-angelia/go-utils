package igorm

import "sync"

// grom pool
var gormPool sync.Map

type Component struct {
	config *config
	locker sync.Mutex
}

// newComponent ...
func newComponent(config *config) *Component {
	return &Component{
		config: config,
	}
}
