package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// logger supports 4 levels. Default level is LevelInformational.
const (
	LevelError = iota
	LevelWarning
	LevelInformational
	LevelDebug
)

// 日志
var Log *Logger

// 预定义字符串
var PreStr string

// Logger ...
type Logger struct {
	// level stores log level.
	level int

	// 4 levels of logger.
	err   *log.Logger
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger

	// depth is the count of the number of frames to skip when computing the file name and line number if Llongfile or Lshortfile is set; a value of 1 will print the details for the caller of Output.
	depth int

	w io.Writer
}

func GetWriter() io.Writer {
	if Log.w == nil {
		return os.Stdout
	}
	return Log.w
}

// NewLogger makes a new logger prints to stdout.
func NewLoggerStdoud(flag int, depth int) *Logger {
	logger := NewWriterLogger(os.Stdout, flag, depth)

	// 设置 writer， 给三方调用
	logger.w = os.Stdout

	return logger
}

// NewFileLogger makes a new file logger, it prints to file lfn. File will auto rotate by size maxsize.
// maxsize is the maximum size in megabytes of the log file
func NewFileLoggerJack(lfn string, maxsize int, flag int, depth int) *Logger {
	jack := &lumberjack.Logger{
		Filename: lfn,
		MaxSize:  maxsize, // megabytes
	}

	logger := NewWriterLogger(jack, flag, depth)

	// 设置 writer， 给三方调用
	logger.w = jack

	return logger
}

func NewFileLoggerRotate(path string, flag int, depth int) *Logger {
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(24*7)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	// 普通日志
	log.SetOutput(writer)

	logger := NewWriterLogger(writer, flag, depth)

	// 设置 writer， 给三方调用
	logger.w = writer

	return logger
}

func NewLogger(projectName string, isOn bool) *Logger {
	if Log != nil {
		return Log
	}
	flag := log.LstdFlags | log.Lshortfile | log.Lmsgprefix

	// 设置 flag
	log.SetFlags(flag)

	if isOn {
		//  初始化日志
		// 下面配置日志每隔 1 天轮转一个新文件，保留最近 7 天的日志文件，多余的自动清理掉。
		path := "./" + projectName + ".log"
		Log = NewFileLoggerRotate(path, flag, 2)

		log.Println("记录日志", path)
	} else {
		Log = NewLoggerStdoud(flag, 3)
	}

	return Log
}

// NewWriterLogger makes a new writer file, it prints to writer.
func NewWriterLogger(w io.Writer, flag int, depth int) *Logger {
	logger := new(Logger)
	logger.depth = depth
	if logger.depth <= 0 {
		logger.depth = 2
	}

	logger.err = log.New(w, "[E] ", flag)
	logger.warn = log.New(w, "[W] ", flag)
	logger.info = log.New(w, "[I] ", flag)
	logger.debug = log.New(w, "[D] ", flag)

	logger.SetLevel(LevelInformational)

	return logger
}

// SetLevel sets the log level.
func (ll *Logger) SetLevel(l int) int {
	ll.level = l
	return ll.level
}

// GetLevel gets the current log level name.
func (ll *Logger) GetLevel() string {
	switch ll.level {
	case LevelDebug:
		return "Debug"
	case LevelError:
		return "Error"
	case LevelInformational:
		return "Info"
	case LevelWarning:
		return "Warn"
	}
	return ""
}

// SetPrefix set the logger prefix. Default prefix is [D] for Debug, [I] for Info, [W] for Warn and [E] for Error.
func (ll *Logger) SetPrefix(prefix string) *Logger {
	ll.err.SetPrefix(prefix)
	ll.warn.SetPrefix(prefix)
	ll.info.SetPrefix(prefix)
	ll.debug.SetPrefix(prefix)
	return ll
}

func (ll *Logger) SetPreStr(uid int32, desc string) *Logger {
	ll.err.SetPrefix(fmt.Sprintf("%suid:%d, desc:%s ", "[E] ", uid, desc))
	ll.warn.SetPrefix(fmt.Sprintf("%suid:%d, desc:%s ", "[W] ", uid, desc))
	ll.info.SetPrefix(fmt.Sprintf("%suid:%d, desc:%s ", "[I] ", uid, desc))
	ll.debug.SetPrefix(fmt.Sprintf("%suid:%d, desc:%s ", "[D] ", uid, desc))
	return ll
}

// Error print log with level Error.
func (ll *Logger) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	ll.err.Output(ll.depth, fmt.Sprintf(format, v...))
}

// Warn print log with level Warn.
func (ll *Logger) Warn(format string, v ...interface{}) {
	if LevelWarning > ll.level {
		return
	}
	ll.warn.Output(ll.depth, fmt.Sprintf(format, v...))
}

// Info print log with level Info.
func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelInformational > ll.level {
		return
	}
	ll.info.Output(ll.depth, fmt.Sprintf(format, v...))
}

// Debug print log with level Debug.
func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	ll.debug.Output(ll.depth, fmt.Sprintf(format, v...))
}

// SetJack makes logger writes to new file lfn.
func (ll *Logger) SetJack(lfn string, maxsize int) {
	jack := &lumberjack.Logger{
		Filename: lfn,
		MaxSize:  maxsize, // megabytes
	}

	ll.err.SetOutput(jack)
	ll.warn.SetOutput(jack)
	ll.info.SetOutput(jack)
	ll.debug.SetOutput(jack)
}

// SetFlag sets log flags. For more information, see the sdk https://golang.org/pkg/log/#pkg-constants.
func (ll *Logger) SetFlag(flag int) {
	ll.err.SetFlags(flag)
	ll.warn.SetFlags(flag)
	ll.info.SetFlags(flag)
	ll.debug.SetFlags(flag)
}

/*
// StdLogger is a predefined logger prints to stdout.
var (
	StdLogger = NewLogger(log.LstdFlags|log.Lshortfile, 3)
)

// SetJack sets the StdLogger's writer to file lfn.
func SetJack(lfn string, maxsize int) {
	StdLogger.SetJack(lfn, maxsize)
}

// Errorf print log with level Error.
func Errorf(format string, v ...interface{}) {
	StdLogger.Error(format, v...)
}

// Warnf print log with level Warn.
func Warnf(format string, v ...interface{}) {
	StdLogger.Warn(format, v...)
}

// Infof print log with level Info.
func Infof(format string, v ...interface{}) {
	StdLogger.Info(format, v...)
}

// Debugf print log with level Debug.
func Debugf(format string, v ...interface{}) {
	StdLogger.Debug(format, v...)
}

// Error print log with level Error.
func Error(v ...interface{}) {
	StdLogger.Error(GenerateFmtStr(len(v)), v...)
}

// Warn print log with level Warn.
func Warn(v ...interface{}) {
	StdLogger.Warn(GenerateFmtStr(len(v)), v...)
}

// Info print log with level Info.
func Info(v ...interface{}) {
	StdLogger.Info(GenerateFmtStr(len(v)), v...)
}

// Debug print log with level Debug.
func Debug(v ...interface{}) {
	StdLogger.Debug(GenerateFmtStr(len(v)), v...)
}


// LogLevel sets the log level.
func LogLevel(logLevel string) string {
	logLevel = strings.ToLower(logLevel)
	if len(logLevel) == 0 {
		logLevel = "info"
	}
	updateLevel(logLevel)
	Warn("Set Log Level as", logLevel)
	return logLevel
}

func updateLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "debug":
		StdLogger.SetLevel(LevelDebug)
	case "info":
		StdLogger.SetLevel(LevelInformational)
	case "warn":
		StdLogger.SetLevel(LevelWarning)
	case "error":
		StdLogger.SetLevel(LevelError)
	default:
		StdLogger.SetLevel(LevelInformational)
	}
}
*/
// GenerateFmtStr is a helper function to construct formatter string.
func GenerateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
