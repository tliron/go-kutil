package logging

import (
	"fmt"
	"sync"
)

//
// PartialLogger
//

type PartialLogger interface {
	Critical(message string)
	Error(message string)
	Warning(message string)
	Notice(message string)
	Info(message string)
	Debug(message string)
}

//
// Logger
//

type Logger interface {
	PartialLogger

	Criticalf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Warningf(format string, values ...interface{})
	Noticef(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Debugf(format string, values ...interface{})
}

//
// FullLogger
//

type FullLogger struct {
	Logger PartialLogger
}

// logging.Logger interface

func (self FullLogger) Critical(message string) {
	self.Logger.Critical(message)
}

func (self FullLogger) Criticalf(format string, values ...interface{}) {
	self.Logger.Critical(fmt.Sprintf(format, values...))
}

func (self FullLogger) Error(message string) {
	self.Logger.Error(message)
}

func (self FullLogger) Errorf(format string, values ...interface{}) {
	self.Logger.Error(fmt.Sprintf(format, values...))
}

func (self FullLogger) Warning(message string) {
	self.Logger.Warning(message)
}

func (self FullLogger) Warningf(format string, values ...interface{}) {
	self.Logger.Warning(fmt.Sprintf(format, values...))
}

func (self FullLogger) Notice(message string) {
	self.Logger.Notice(message)
}

func (self FullLogger) Noticef(format string, values ...interface{}) {
	self.Logger.Notice(fmt.Sprintf(format, values...))
}

func (self FullLogger) Info(message string) {
	self.Logger.Info(message)
}

func (self FullLogger) Infof(format string, values ...interface{}) {
	self.Logger.Info(fmt.Sprintf(format, values...))
}

func (self FullLogger) Debug(message string) {
	self.Logger.Debug(message)
}

func (self FullLogger) Debugf(format string, values ...interface{}) {
	self.Logger.Debug(fmt.Sprintf(format, values...))
}

//
// LazyLogger
//

type LazyLogger struct {
	Name   string
	Logger Logger

	once sync.Once
}

func (self *LazyLogger) validate() {
	self.once.Do(func() {
		if self.Logger == nil {
			if backend == nil {
				panic("logging not configured")
			}
			self.Logger = backend.GetLogger(self.Name)
		}
	})
}

// logging.Logger interface

func (self *LazyLogger) Critical(message string) {
	self.validate()
	self.Logger.Critical(message)
}

func (self *LazyLogger) Criticalf(format string, values ...interface{}) {
	self.validate()
	self.Logger.Criticalf(format, values...)
}

func (self *LazyLogger) Error(message string) {
	self.validate()
	self.Logger.Error(message)
}

func (self *LazyLogger) Errorf(format string, values ...interface{}) {
	self.validate()
	self.Logger.Errorf(format, values...)
}

func (self *LazyLogger) Warning(message string) {
	self.validate()
	self.Logger.Warning(message)
}

func (self *LazyLogger) Warningf(format string, values ...interface{}) {
	self.validate()
	self.Logger.Warningf(format, values...)
}

func (self *LazyLogger) Notice(message string) {
	self.validate()
	self.Logger.Notice(message)
}

func (self *LazyLogger) Noticef(format string, values ...interface{}) {
	self.validate()
	self.Logger.Noticef(format, values...)
}

func (self *LazyLogger) Info(message string) {
	self.validate()
	self.Logger.Info(message)
}

func (self *LazyLogger) Infof(format string, values ...interface{}) {
	self.validate()
	self.Logger.Infof(format, values...)
}

func (self *LazyLogger) Debug(message string) {
	self.validate()
	self.Logger.Debug(message)
}

func (self *LazyLogger) Debugf(format string, values ...interface{}) {
	self.validate()
	self.Logger.Debugf(format, values...)
}
