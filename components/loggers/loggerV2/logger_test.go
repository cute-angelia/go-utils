package loggerV2

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// online false
	loggerObj := NewLogger(WithProject("test-local"), WithIsOnline(false))
	loggerObj.Info("start")
	loggerObj.WithTime(time.Now()).WithField("title", "for test").Error("not sup")

	// online true
	logger2 := NewLogger(WithProject("test-local"), WithIsOnline(true))
	for i := 0; i < 100; i++ {
		logger2.Info("start")
		logger2.WithTime(time.Now()).WithField("title", "for test").Error("not sup")
	}
}
