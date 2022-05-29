package iruntime

import (
	"runtime"
	"strings"
)

// IsWindows windows系统
func IsWindows() bool {
	sysType := runtime.GOOS
	if strings.Contains(sysType, "windows") {
		return true
	}
	return false
}

func IsLinux() bool {
	sysType := runtime.GOOS
	if strings.Contains(sysType, "linux") {
		return true
	}
	return false
}
