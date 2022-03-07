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
)

var logOnce sync.Once
var logger *zerolog.Logger

type Component struct {
	config *config
}

// GetLogger 开放方法
func GetLogger() *zerolog.Logger {
	if logger == nil {
		log.Println("日志未初始化，将启动`stdout`输出模式")
		logger = newComponent(DefaultConfig()).NewLogger()
	}
	return logger
}

func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config
	return comp
}

func (self *Component) NewLogger() *zerolog.Logger {
	logOnce.Do(func() {
		// 开发模式自定义日志格式
		output := self.formatLogger(os.Stdout)
		// 原生日志支持
		log.SetFlags(0)
		log.SetOutput(zerolog.New(output).With().Timestamp().CallerWithSkipFrameCount(4).Logger())

		var writers []io.Writer
		if self.config.isOnline {

			if self.config.fileJson {
				// 线上 json 模式
				writers = append(writers, self.newRollingFile())
			} else {
				// 线上输出模式
				writers = append(writers, self.formatLogger(self.newRollingFile()))
			}
		} else {
			writers = append(writers, output)
		}
		mw := zerolog.MultiLevelWriter(writers...)

		// 配置
		zerolog.SetGlobalLevel(self.config.level)
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
		zerolog.CallerMarshalFunc = func(file string, line int) string {
			var buffer bytes.Buffer
			buffer.WriteString(path.Base(file))
			buffer.WriteString(":")
			buffer.WriteString(strconv.Itoa(line))
			return buffer.String()
		}
		ilog := zerolog.New(mw).With().Timestamp().Caller().Logger()
		logger = &ilog
	})
	return logger
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

func (self *Component) newRollingFile() io.Writer {
	return &lumberjack.Logger{
		Filename:   "./log_" + self.config.project + ".log",
		MaxBackups: self.config.maxBackups, // files
		MaxSize:    self.config.maxSize,    // megabytes
		MaxAge:     self.config.maxAge,     // days
	}
}
