package logger2

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
)

type logger struct {
	*logrus.Logger
	project    string
	isOnline   bool
	maxSize    int
	maxBackups int // log nums
	maxAge     int //days
}

func NewLogger(options ...func(*logger)) *logger {
	ilogger := &logger{}
	for _, o := range options {
		o(ilogger)
	}

	if ilogger.maxSize <= 0 {
		ilogger.maxSize = 100
	}

	if ilogger.maxBackups <= 0 {
		ilogger.maxBackups = 10
	}

	if ilogger.maxAge <= 0 {
		ilogger.maxAge = 10
	}

	l := logrus.New()
	l.ReportCaller = true

	if !ilogger.isOnline {
		l.Out = os.Stdout
		l.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	} else {
		l.Out = &lumberjack.Logger{
			Filename:   fmt.Sprintf("./%s.log", ilogger.project),
			MaxSize:    ilogger.maxSize, // megabytes
			MaxBackups: 10,
			MaxAge:     10, //days
			LocalTime:  true,
		}
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	}

	// logger
	ilogger.Logger = l

	return ilogger
}

func WithProject(project string) func(*logger) {
	return func(s *logger) {
		s.project = project
	}
}

func WithIsOnline(isOnline bool) func(*logger) {
	return func(s *logger) {
		s.isOnline = isOnline
	}
}

func WithMaxSize(maxSize int) func(*logger) {
	return func(s *logger) {
		s.maxSize = maxSize
	}
}
func WithMaxBackups(maxBackups int) func(*logger) {
	return func(s *logger) {
		s.maxBackups = maxBackups
	}
}
func WithMaxAge(maxAge int) func(*logger) {
	return func(s *logger) {
		s.maxAge = maxAge
	}
}
