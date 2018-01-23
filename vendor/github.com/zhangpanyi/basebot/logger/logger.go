package logger

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/rifflock/lfshook"
	"github.com/senlinms/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	PanicLevel = uint8(logrus.PanicLevel)
	FatalLevel = uint8(logrus.FatalLevel)
	ErrorLevel = uint8(logrus.ErrorLevel)
	WarnLevel  = uint8(logrus.WarnLevel)
	InfoLevel  = uint8(logrus.InfoLevel)
	DebugLevel = uint8(logrus.DebugLevel)
)

// CreateLoggerOnce 创建日志记录器
func CreateLoggerOnce(level, filelevel uint8) {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	base := filepath.Base(os.Args[0])
	name := strings.TrimSuffix(base, filepath.Ext(base))
	path := "logs/" + name + "/" + tm.Format("20060102150405") + "/"

	once.Do(func() {
		globalLogger = &logger{
			file:      newLoggerOfLFShook(1048576000, 100, 365, path),
			console:   newLoggerOfConsole(),
			fileLevel: log.Level(level),
		}
		globalLogger.console.SetLevel(log.Level(level))
		globalLogger.file.SetLevel(log.Level(filelevel))
	})
}

// Debug 输出Debug日志
func Debug(v ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.DebugLevel).Debug(v)
	}
}

// Debugf 格式化输出Debug日志
func Debugf(format string, params ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.DebugLevel).Debugf(format, params...)
	}
}

// Info 输出Info日志
func Info(v ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.InfoLevel).Info(v)
	}
}

// Infof 格式化输出Info日志
func Infof(format string, params ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.InfoLevel).Infof(format, params...)
	}
}

// Warn 输出Warn日志
func Warn(v ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.WarnLevel).Warn(v)
	}
}

// Warnf 格式化输出Warn日志
func Warnf(format string, params ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.WarnLevel).Warnf(format, params...)
	}
}

// Error 输出Error日志
func Error(v ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.ErrorLevel).Error(v)
	}
}

// Errorf 格式化输出Error日志
func Errorf(format string, params ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.ErrorLevel).Errorf(format, params...)
	}
}

// Fatal 输出Fatal日志
func Fatal(v ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.FatalLevel).Fatal(v)
	}
}

// Fatalf 格式化输出Fatal日志
func Fatalf(format string, params ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.FatalLevel).Fatalf(format, params...)
	}
}

// Panic 输出Panic日志
func Panic(v ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.PanicLevel).Panic(v)
	}
}

// Panicf 格式化输出Panic日志
func Panicf(format string, params ...interface{}) {
	if globalLogger != nil {
		globalLogger.available(log.PanicLevel).Panicf(format, params...)
	}
}

// 日志选项
type logger struct {
	file      *log.Logger // 文件日志
	console   *log.Logger // 控制台日志
	fileLevel log.Level   // 文件日志级别
}

// 获取可用的日志记录器
func (lg *logger) available(level log.Level) *log.Logger {
	if level >= lg.fileLevel {
		return lg.file
	}
	return lg.console
}

var once sync.Once
var globalLogger *logger

// 创建终端记录器
func newLoggerOfConsole() *log.Logger {
	lg := log.New()
	for _, level := range log.AllLevels {
		lg.Level |= level
	}
	return lg
}

// 创建文件记录器
func newLoggerOfLFShook(maxsize int, maxbackup int, maxage int, path string) *log.Logger {
	lg := log.New()
	writerMap := lfshook.WriterMap{}
	for _, level := range log.AllLevels {
		lg.Level |= level
		writer := &lumberjack.Logger{
			Filename:   path + level.String() + ".log",
			MaxSize:    maxsize,
			MaxBackups: maxbackup,
			MaxAge:     maxage,
		}
		writerMap[level] = writer
	}
	lg.Formatter = &log.JSONFormatter{}
	lg.Hooks.Add(lfshook.NewHook(writerMap, nil))
	return lg
}
