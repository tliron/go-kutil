package logging

import (
	"fmt"
	"sync"
)

var backend Backend

func SetBackend(backend_ Backend) {
	backend = backend_
}

func Configure(verbosity int, path *string) {
	backend.Configure(verbosity, path)
}

func SetMaxLevel(name string, level Level) {
	backend.SetMaxLevel(name, level)
}

func GetLogger(name string) Logger {
	return &LazyLogger{
		Name: name,
	}
}

func GetLoggerf(format string, arguments ...interface{}) Logger {
	return GetLogger(fmt.Sprintf(format, arguments...))
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
