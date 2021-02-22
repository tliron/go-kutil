package gologging

import (
	loggingpkg "github.com/op/go-logging"
)

//
// Logger
//

type Logger struct {
	Logger *loggingpkg.Logger
}

// logging.Logger interface

func (self Logger) Critical(message string) {
	self.Logger.Critical(message)
}

func (self Logger) Criticalf(format string, values ...interface{}) {
	self.Logger.Criticalf(format, values...)
}

func (self Logger) Error(message string) {
	self.Logger.Error(message)
}

func (self Logger) Errorf(format string, values ...interface{}) {
	self.Logger.Errorf(format, values...)
}

func (self Logger) Warning(message string) {
	self.Logger.Warning(message)
}

func (self Logger) Warningf(format string, values ...interface{}) {
	self.Logger.Warningf(format, values...)
}

func (self Logger) Notice(message string) {
	self.Logger.Notice(message)
}

func (self Logger) Noticef(format string, values ...interface{}) {
	self.Logger.Noticef(format, values...)
}

func (self Logger) Info(message string) {
	self.Logger.Info(message)
}

func (self Logger) Infof(format string, values ...interface{}) {
	self.Logger.Infof(format, values...)
}

func (self Logger) Debug(message string) {
	self.Logger.Debug(message)
}

func (self Logger) Debugf(format string, values ...interface{}) {
	self.Logger.Debugf(format, values...)
}
