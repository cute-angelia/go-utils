package helper

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