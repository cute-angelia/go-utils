package itime

import (
	"errors"
	"time"
)

// errTimeout indicates a timeout.
var errTimeout = errors.New("timeout")

type (
	// Ticker interface wraps the Chan and Stop methods.
	Ticker interface {
		Chan() <-chan time.Time
		Stop()
	}

	// FakeTicker interface is used for unit testing.
	FakeTicker interface {
		Ticker
		Done()
		Tick()
		Wait(d time.Duration) error
	}

	fakeTicker struct {
		c    chan time.Time
		done chan struct{}
	}

	realTicker struct {
		*time.Ticker
	}
)

// NewTicker returns a Ticker.
func NewTicker(d time.Duration) Ticker {
	return &realTicker{
		Ticker: time.NewTicker(d),
	}
}

func (rt *realTicker) Chan() <-chan time.Time {
	return rt.C
}

// NewFakeTicker returns a FakeTicker.
func NewFakeTicker() FakeTicker {
	return &fakeTicker{
		c:    make(chan time.Time, 1),
		done: make(chan struct{}, 1),
	}
}

func (ft *fakeTicker) Chan() <-chan time.Time {
	return ft.c
}

func (ft *fakeTicker) Done() {
	ft.done <- Placeholder
}

func (ft *fakeTicker) Stop() {
	close(ft.c)
}

func (ft *fakeTicker) Tick() {
	ft.c <- time.Now()
}

func (ft *fakeTicker) Wait(d time.Duration) error {
	select {
	case <-time.After(d):
		return errTimeout
	case <-ft.done:
		return nil
	}
}
