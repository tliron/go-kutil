package logging

import (
	"fmt"
)

const DEFAULT_SUBLOGGER_FORMAT = "{%s} %s"

//
// SubLogger
//

type SubLogger struct {
	logger Logger
	domain string
}

func NewSubLogger(logger Logger, domain string) SubLogger {
	if subLogger, ok := logger.(SubLogger); ok {
		domain = subLogger.domain + "." + domain
		logger = subLogger.logger
	}

	return SubLogger{
		logger: logger,
		domain: domain,
	}
}

// Logger interface

func (self SubLogger) NewMessage(level Level, depth int) Message {
	if message := self.logger.NewMessage(level, depth); message != nil {
		message.Set("domain", self.domain)
		return message
	} else {
		return nil
	}
}

func (self SubLogger) Log(level Level, depth int, message string) {
	if message_ := self.NewMessage(level, depth+1); message_ != nil {
		message_.Set("message", message)
		message_.Send()
	}
}

func (self SubLogger) Logf(level Level, depth int, format string, values ...interface{}) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("message", fmt.Sprintf(format, values...))
		message.Send()
	}
}

func (self SubLogger) Critical(message string) {
	self.Log(Critical, 1, message)
}

func (self SubLogger) Criticalf(format string, values ...interface{}) {
	self.Logf(Critical, 1, format, values...)
}

func (self SubLogger) Error(message string) {
	self.Log(Error, 1, message)
}

func (self SubLogger) Errorf(format string, values ...interface{}) {
	self.Logf(Error, 1, format, values...)
}

func (self SubLogger) Warning(message string) {
	self.Log(Warning, 1, message)
}

func (self SubLogger) Warningf(format string, values ...interface{}) {
	self.Logf(Warning, 1, format, values...)
}

func (self SubLogger) Notice(message string) {
	self.Log(Notice, 1, message)
}

func (self SubLogger) Noticef(format string, values ...interface{}) {
	self.Logf(Notice, 1, format, values...)
}

func (self SubLogger) Info(message string) {
	self.Log(Info, 1, message)
}

func (self SubLogger) Infof(format string, values ...interface{}) {
	self.Logf(Info, 1, format, values...)
}

func (self SubLogger) Debug(message string) {
	self.Log(Debug, 1, message)
}

func (self SubLogger) Debugf(format string, values ...interface{}) {
	self.Logf(Debug, 1, format, values...)
}

func (self SubLogger) AllowLevel(level Level) bool {
	return self.logger.AllowLevel(level)
}

func (self SubLogger) SetMaxLevel(level Level) {
	self.logger.SetMaxLevel(level)
}
