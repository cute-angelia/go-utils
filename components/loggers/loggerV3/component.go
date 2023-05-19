package loggerV3

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var component *Component

var logOnce sync.Once
var loggerError *zerolog.Logger

type Component struct {
	config *config
	logger *zerolog.Logger
}

// GetLogger 开放方法
func GetLogger() *zerolog.Logger {
	if component == nil || component.logger == nil {
		log.Println("日志未初始化，启动默认输出保存为文件log_default.log")
		component = newComponent(DefaultConfig())
	}

	return component.logger
}

func newComponent(config *config) *Component {
	cpt := &Component{}
	cpt.config = config
	cpt.logger = cpt.newLogger(config)

	component = cpt
	return component
}

func (self *Component) newLogger(config *config) *zerolog.Logger {
	logOnce.Do(func() {
		ilog := self.makeMainLogger("/log_" + config.Project + ".log")
		// hook error
		if config.HookError {
			ilog = ilog.Hook(ErrorHook{})
			loge := self.makeErrorLogger("/error/log_error_" + config.Project + ".log")
			loggerError = &loge
		}
		self.logger = &ilog
	})
	return self.logger
}

// makeMainLoger 主 logger
func (self *Component) makeMainLogger(logName string) zerolog.Logger {
	var writers []io.Writer
	// 原生日志支持
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix | log.Lmicroseconds)
	log.Println(fmt.Sprintf("[%s] 初始化，项目:%s 路径：%s 线上模式:%v FileJson:%v", ComponentName, self.config.Project, self.config.LogPath, self.config.IsOnline, self.config.FileJson))
	if self.config.IsOnline {
		var w io.Writer
		// 线上 json 模式
		if self.config.FileJson {
			w = self.newRollingFile(logName)
			log.SetOutput(w)
		} else {
			// 线上输出模式
			w = self.newRollingFile(logName)
			log.SetOutput(w)
			w = self.formatLogger(w)
		}
		writers = append(writers, w)
	} else {
		// 开发模式自定义日志格式
		output := self.formatLogger(os.Stdout)
		// 关闭 caller
		// log.Println("callerWithSkipFrameCount:", callerWithSkipFrameCount, callerWithSkipFrameCount == -1)
		log.SetOutput(zerolog.New(output).With().Timestamp().CallerWithSkipFrameCount(4).Logger())

		writers = append(writers, output)
	}
	mw := zerolog.MultiLevelWriter(writers...)

	// 配置
	zerolog.SetGlobalLevel(self.config.Level)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		var buffer bytes.Buffer
		buffer.WriteString(path.Base(file))
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(line))
		return buffer.String()
	}
	ilog := zerolog.New(mw).With().Timestamp().Caller().Logger()
	return ilog
}

// makeErrorLogger 需求，error 级别单独文件存储
func (self *Component) makeErrorLogger(logName string) zerolog.Logger {
	var writers []io.Writer
	if self.config.IsOnline {
		var w io.Writer
		// 线上 json 模式
		if self.config.FileJson {
			w = self.newRollingFile(logName)
		} else {
			// 线上输出模式
			w = self.newRollingFile(logName)
			w = self.formatLogger(w)
		}
		writers = append(writers, w)
	} else {
		// 开发模式自定义日志格式
		output := self.formatLogger(os.Stdout)
		// 关闭 caller
		log.SetOutput(zerolog.New(output).With().Timestamp().Logger())
		writers = append(writers, output)
	}
	mw := zerolog.MultiLevelWriter(writers...)

	// 配置
	zerolog.SetGlobalLevel(self.config.Level)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		var buffer bytes.Buffer
		buffer.WriteString(path.Base(file))
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(line))
		return buffer.String()
	}
	ilog := zerolog.New(mw).With().Timestamp().Logger()
	return ilog
}

func (self *Component) formatLogger(out io.Writer) io.Writer {
	// 开发模式自定义日志格式
	output := zerolog.ConsoleWriter{Out: out, TimeFormat: "2006-01-02 15:04:05.000"}
	output.FormatLevel = func(i interface{}) string {
		if i == nil {
			i = "info"
		}
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	return output
}

func (self *Component) newRollingFile(logName string) io.Writer {
	ljLogger := lumberjack.Logger{
		Filename:   self.config.LogPath + "/" + logName,
		MaxBackups: self.config.MaxBackups, // files
		MaxSize:    self.config.MaxSize,    // megabytes
		MaxAge:     self.config.MaxAge,     // days
	}

	// 每日分割
	if self.config.Everyday {
		go func() {
			for {
				now := time.Now()
				// 计算下一个零点
				next := now.Add(time.Hour * 24)
				next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
				t := time.NewTimer(next.Sub(now))
				<-t.C
				ljLogger.Rotate()
				t.Stop()
			}
		}()
	}

	return &ljLogger
}
