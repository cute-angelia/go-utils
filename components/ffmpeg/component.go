package ffmpeg

import (
	"github.com/gotomicro/ego/core/elog"
	"sync"
)

type Component struct {
	config *config
	logger *elog.Component
	locker sync.Mutex
}

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	return &Component{
		config: config,
		logger: logger,
	}
}

func (c Component) GetFfmpegPath() string {
	return c.config.FfmpegPath
}
