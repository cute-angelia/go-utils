package ibadger

import (
	"github.com/dgraph-io/badger/v3"
	"time"
)

// config options
type config struct {
	Path string

	// Options are the badger datastore options, reexported here for convenience.
	// Please refer to the Badger docs to see what this is for
	GcDiscardRatio float64

	// Interval between GC cycles
	//
	// If zero, the datastore will perform no automatic garbage collection.
	GcInterval time.Duration

	// Sleep time between rounds of a single GC cycle.
	//
	// If zero, the datastore will only perform one round of GC per
	// GcInterval.
	GcSleep time.Duration

	badger.Options
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	DefaultOptions := config{
		Path:           "./badger",
		GcDiscardRatio: 0.2,
		GcInterval:     15 * time.Minute,
		GcSleep:        10 * time.Second,
		Options:        badger.LSMOnlyOptions(""),
	}

	// path
	DefaultOptions.Options.Dir = DefaultOptions.Path
	DefaultOptions.Options.ValueDir = DefaultOptions.Path

	// This is to optimize the database on close so it can be opened
	// read-only and efficiently queried. We don't do that and hanging on
	// stop isn't nice.
	// DefaultOptions.Options.CompactL0OnClose = false

	return &DefaultOptions
}
