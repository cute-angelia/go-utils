package ffmpeg

import (
	"sync"
)

type Component struct {
	config *config
	locker sync.Mutex
}

// newComponent ...
func newComponent(compName string, config *config) *Component {
	return &Component{
		config: config,
	}
}

func (c Component) GetFfmpegPath() string {
	return c.config.FfmpegPath
}
