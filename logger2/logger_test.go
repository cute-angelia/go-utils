package logger2

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// online false
	logger := NewLogger(WithProject("test-local"), WithIsOnline(false))
	logger.Info("start")
	logger.WithTime(time.Now()).WithField("title", "for test").Error("not sup")

	// online true
	logger2 := NewLogger(WithProject("test-local"), WithIsOnline(true))
	for i := 0; i < 100; i++ {
		logger2.Info("start")
		logger2.WithTime(time.Now()).WithField("title", "for test").Error("not sup")
	}
}

