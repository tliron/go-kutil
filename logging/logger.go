package logging

import "fmt"

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
