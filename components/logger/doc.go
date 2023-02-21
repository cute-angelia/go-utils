/*
A logger package that compatible with sdk log package. It supports level logging and file auto rotating.
It also has a predefined 'standard' Logger called StdLogger accessible through helper functions Debug[f], Info[f], Warn[f], Error[f], LogLevel and SetJack.
It supports 4 level:
	LevelError = iota
	LevelWarning
	LevelInformational
	LevelDebug
You can use LogLevel to handle the log level.
File rotating based on package gopkg.in/natefinch/lumberjack.v2, you can control file settings by using SetJack.
Quick start
Use StdLogger:
	import "github.com/mnhkahn/gogogo/logger"
	logger.Info("hello, world.")
Defined our own logger:
	l := logger.NewWriterLogger(w, flag, 3)
	l.Info("hello, world")
Log flags compatible with sdk log package, https://golang.org/pkg/log/#pkg-constants.
See also
For more information, goto godoc https://godoc.org/github.com/mnhkahn/gogogo/logger
Chinese details, goto http://blog.cyeam.com//golang/2017/07/14/go-log?utm_source=github


import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"os"
	"time"
)

func SetLogWithOutputRotateLog(path string) *rotatelogs.RotateLogs {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(24*7)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	return writer
}

func SetLogWithOsOut() *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return os.Stdout
}
*/
package logger
