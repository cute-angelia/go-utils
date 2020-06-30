package helper

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"time"
)

func SetLog(path string) {
	if len(path) == 0 {
		path = "./go.log"
	}
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(24*7)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	log.SetOutput(writer)
}
