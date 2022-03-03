package nuts

import (
	"github.com/xujiajun/nutsdb"
)

// PackageName ..
const PackageName = "contrib.nuts"

// config
type config struct {
	Key     string
	Options nutsdb.Options
}

// DefaultConfig ...
func DefaultConfig() *config {
	c := config{
		Key:     "default",
		Options: nutsdb.DefaultOptions,
	}
	return &c
}
