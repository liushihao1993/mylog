package mylog

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type OptionApplier interface {
	Apply(defaultLogger *Logger)
}
type LoggerOption struct {
	LogName Name
	Logger  *lumberjack.Logger
}

func (option *LoggerOption) Apply(defaultLogger *Logger) {
	defaultInfoLogger := defaultLogger.defaultInfoLogger
	errLogger := defaultLogger.defaultErrorLogger
	prefix := fmt.Sprintf("[%s] ", option.LogName)
	defaultLogger.logMap[option.LogName] = &loggerWithLevel{
		infoLogger:  &logger{Writer: io.MultiWriter(option.Logger, defaultInfoLogger), prefix: prefix},
		warnLogger:  &logger{Writer: io.MultiWriter(option.Logger, errLogger), prefix: prefix},
		errorLogger: &logger{Writer: io.MultiWriter(option.Logger, errLogger), prefix: prefix},
	}
}

type FileLinOption struct {
	HideFileLine bool
}

func (f FileLinOption) Apply(defaultLogger *Logger) {
	defaultLogger.hideFileLine.Store(f.HideFileLine)
}

type FunctionOption struct {
	HideFunction bool
}

func (f FunctionOption) Apply(defaultLogger *Logger) {
	defaultLogger.hideFunction.Store(f.HideFunction)
}
