package logger

import "testing"

const Project = "TEST"

// go test -run TestNew
func TestNew(t *testing.T) {
	l := NewLogger(Project, false)
	l.SetPreStr(123456, ":for test:")
	l.Info("%s", "new logger")
}