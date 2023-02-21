package ifile

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"runtime"
)

func GetHomeDir() string {
	if homedirstr, err := homedir.Dir(); err != nil {
		return "./"
	} else {
		return homedirstr
	}
}

// 获取用户文件夹
func GetUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}
