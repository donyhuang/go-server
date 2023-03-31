package log

import (
	"context"
	"github.com/donyhuang/go-server/pkg/reflectx"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

const (
	GinLogEntryWithTraceKey = "LogEntryWithTrace"
)

type logEntryContextKeyType int

const logEntryWithTraceKey logEntryContextKeyType = iota

var (
	once           sync.Once
	Logger         *logrus.Logger
	LoggerRecovery *logrus.Logger
	LoggerEntry    *logrus.Entry
	defaultConf    = Conf{
		Path:       "./logs/server.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     3,
		Level:      "trace",
	}
)

type Conf struct {
	Path       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      string
}
type Option func(*Conf)

func GetLogRecovery() io.Writer {
	once.Do(func() {
		out := &lumberjack.Logger{
			Filename:   "./logs/gin_recovery.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     3,
		}
		LoggerRecovery = &logrus.Logger{
			Out:       out,
			Formatter: &logrus.TextFormatter{},
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		}
	})
	return LoggerRecovery.Out
}

func InitLog(c Conf) {
	l := defaultConf
	_ = reflectx.SetStructNotEmpty(&l, c)
	level, _ := logrus.ParseLevel(l.Level)
	out := &lumberjack.Logger{
		Filename:   l.Path,
		MaxSize:    l.MaxSize,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
	}

	Logger = &logrus.Logger{
		Out:       out,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
	LoggerEntry = logrus.NewEntry(Logger)
}

func GetLogEntryFromContext(c context.Context) *logrus.Entry {
	switch c.(type) {
	case *gin.Context:
		p, ok := c.(*gin.Context).Get(GinLogEntryWithTraceKey)
		if ok {
			logEntry, ok := p.(*logrus.Entry)
			if ok {
				return logEntry
			}
		}
		return LoggerEntry
	default:
		logEntry, ok := c.Value(logEntryWithTraceKey).(*logrus.Entry)
		if ok {
			return logEntry
		}
		if LoggerEntry == nil {
			return &logrus.Entry{Logger: &logrus.Logger{Out: os.Stdout}}
		}
		return LoggerEntry
	}
}
func InjectLogEntryToContext(traceId string, c context.Context) context.Context {
	entry := LoggerEntry.WithField("TraceId", traceId)
	return context.WithValue(c, logEntryWithTraceKey, entry)
}
func InjectLogEntryToGinContext(traceId string, c *gin.Context) {
	c.Set(GinLogEntryWithTraceKey, LoggerEntry.WithField("TraceId", traceId))
}
