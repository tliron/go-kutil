package logging

import (
	"fmt"

	loggingpkg "github.com/op/go-logging"
)

//
// Logger
//

type Logger interface {
	Critical(message string)
	Criticalf(format string, values ...interface{})
	Error(message string)
	Errorf(format string, values ...interface{})
	Warning(message string)
	Warningf(format string, values ...interface{})
	Notice(message string)
	Noticef(format string, values ...interface{})
	Info(message string)
	Infof(format string, values ...interface{})
	Debug(message string)
	Debugf(format string, values ...interface{})
}

func GetLogger(name string) Logger {
	return GoLoggingLogger{loggingpkg.MustGetLogger(name)}
}

func GetLoggerf(format string, arguments ...interface{}) Logger {
	return GetLogger(fmt.Sprintf(format, arguments...))
}

//
// GoLoggingLogger
//

type GoLoggingLogger struct {
	logger *loggingpkg.Logger
}

// Logger interface

func (self GoLoggingLogger) Critical(message string) {
	self.logger.Critical(message)
}

func (self GoLoggingLogger) Criticalf(format string, values ...interface{}) {
	self.logger.Criticalf(format, values...)
}

func (self GoLoggingLogger) Error(message string) {
	self.logger.Error(message)
}

func (self GoLoggingLogger) Errorf(format string, values ...interface{}) {
	self.logger.Errorf(format, values...)
}

func (self GoLoggingLogger) Warning(message string) {
	self.logger.Warning(message)
}

func (self GoLoggingLogger) Warningf(format string, values ...interface{}) {
	self.logger.Warningf(format, values...)
}

func (self GoLoggingLogger) Notice(message string) {
	self.logger.Notice(message)
}

func (self GoLoggingLogger) Noticef(format string, values ...interface{}) {
	self.logger.Noticef(format, values...)
}

func (self GoLoggingLogger) Info(message string) {
	self.logger.Info(message)
}

func (self GoLoggingLogger) Infof(format string, values ...interface{}) {
	self.logger.Infof(format, values...)
}

func (self GoLoggingLogger) Debug(message string) {
	self.logger.Debug(message)
}

func (self GoLoggingLogger) Debugf(format string, values ...interface{}) {
	self.logger.Debugf(format, values...)
}
