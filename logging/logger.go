package logging

import "fmt"

//
// Logger
//

type Logger interface {
	NewMessage(level Level, depth int) Message
	AllowLevel(level Level) bool
	SetMaxLevel(level Level)

	// For unstructured logging

	Log(level Level, depth int, message string)
	Logf(level Level, depth int, format string, values ...interface{})

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

//
// BackendLogger
//

type BackendLogger struct {
	name []string
}

func NewBackendLogger(id []string) BackendLogger {
	return BackendLogger{name: id}
}

// Logger interface

func (self BackendLogger) NewMessage(level Level, depth int) Message {
	return NewMessage(self.name, level, depth)
}

func (self BackendLogger) AllowLevel(level Level) bool {
	return AllowLevel(self.name, level)
}

func (self BackendLogger) SetMaxLevel(level Level) {
	SetMaxLevel(self.name, level)
}

func (self BackendLogger) Log(level Level, depth int, message string) {
	if message_ := self.NewMessage(level, depth+1); message_ != nil {
		message_.Set("message", message)
		message_.Send()
	}
}

func (self BackendLogger) Logf(level Level, depth int, format string, values ...interface{}) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("message", fmt.Sprintf(format, values...))
		message.Send()
	}
}

func (self BackendLogger) Critical(message string) {
	self.Log(Critical, 1, message)
}

func (self BackendLogger) Criticalf(format string, values ...interface{}) {
	self.Logf(Critical, 1, format, values...)
}

func (self BackendLogger) Error(message string) {
	self.Log(Error, 1, message)
}

func (self BackendLogger) Errorf(format string, values ...interface{}) {
	self.Logf(Error, 1, format, values...)
}

func (self BackendLogger) Warning(message string) {
	self.Log(Warning, 1, message)
}

func (self BackendLogger) Warningf(format string, values ...interface{}) {
	self.Logf(Warning, 1, format, values...)
}

func (self BackendLogger) Notice(message string) {
	self.Log(Notice, 1, message)
}

func (self BackendLogger) Noticef(format string, values ...interface{}) {
	self.Logf(Notice, 1, format, values...)
}

func (self BackendLogger) Info(message string) {
	self.Log(Info, 1, message)
}

func (self BackendLogger) Infof(format string, values ...interface{}) {
	self.Logf(Info, 1, format, values...)
}

func (self BackendLogger) Debug(message string) {
	self.Log(Debug, 1, message)
}

func (self BackendLogger) Debugf(format string, values ...interface{}) {
	self.Logf(Debug, 1, format, values...)
}

//
// ScopeLogger
//

type ScopeLogger struct {
	logger Logger
	scope  string
}

func NewScopeLogger(logger Logger, scope string) ScopeLogger {
	if subLogger, ok := logger.(ScopeLogger); ok {
		scope = subLogger.scope + "." + scope
		logger = subLogger.logger
	}

	return ScopeLogger{
		logger: logger,
		scope:  scope,
	}
}

// Logger interface

func (self ScopeLogger) NewMessage(level Level, depth int) Message {
	if message := self.logger.NewMessage(level, depth); message != nil {
		message.Set("scope", self.scope)
		return message
	} else {
		return nil
	}
}

func (self ScopeLogger) AllowLevel(level Level) bool {
	return self.logger.AllowLevel(level)
}

func (self ScopeLogger) SetMaxLevel(level Level) {
	self.logger.SetMaxLevel(level)
}

func (self ScopeLogger) Log(level Level, depth int, message string) {
	if message_ := self.NewMessage(level, depth+1); message_ != nil {
		message_.Set("message", message)
		message_.Send()
	}
}

func (self ScopeLogger) Logf(level Level, depth int, format string, values ...interface{}) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("message", fmt.Sprintf(format, values...))
		message.Send()
	}
}

func (self ScopeLogger) Critical(message string) {
	self.Log(Critical, 1, message)
}

func (self ScopeLogger) Criticalf(format string, values ...interface{}) {
	self.Logf(Critical, 1, format, values...)
}

func (self ScopeLogger) Error(message string) {
	self.Log(Error, 1, message)
}

func (self ScopeLogger) Errorf(format string, values ...interface{}) {
	self.Logf(Error, 1, format, values...)
}

func (self ScopeLogger) Warning(message string) {
	self.Log(Warning, 1, message)
}

func (self ScopeLogger) Warningf(format string, values ...interface{}) {
	self.Logf(Warning, 1, format, values...)
}

func (self ScopeLogger) Notice(message string) {
	self.Log(Notice, 1, message)
}

func (self ScopeLogger) Noticef(format string, values ...interface{}) {
	self.Logf(Notice, 1, format, values...)
}

func (self ScopeLogger) Info(message string) {
	self.Log(Info, 1, message)
}

func (self ScopeLogger) Infof(format string, values ...interface{}) {
	self.Logf(Info, 1, format, values...)
}

func (self ScopeLogger) Debug(message string) {
	self.Log(Debug, 1, message)
}

func (self ScopeLogger) Debugf(format string, values ...interface{}) {
	self.Logf(Debug, 1, format, values...)
}

//
// MockLogger
//

var MOCK_LOGGER MockLogger

type MockLogger struct{}

// Logger interface

func (self MockLogger) NewMessage(level Level, depth int) Message {
	return nil
}

func (self MockLogger) AllowLevel(level Level) bool {
	return false
}

func (self MockLogger) SetMaxLevel(level Level) {
}

func (self MockLogger) Log(level Level, depth int, message string) {
}

func (self MockLogger) Logf(level Level, depth int, format string, values ...interface{}) {
}

func (self MockLogger) Critical(message string) {
}

func (self MockLogger) Criticalf(format string, values ...interface{}) {
}

func (self MockLogger) Error(message string) {
}

func (self MockLogger) Errorf(format string, values ...interface{}) {
}

func (self MockLogger) Warning(message string) {
}

func (self MockLogger) Warningf(format string, values ...interface{}) {
}

func (self MockLogger) Notice(message string) {
}

func (self MockLogger) Noticef(format string, values ...interface{}) {
}

func (self MockLogger) Info(message string) {
}

func (self MockLogger) Infof(format string, values ...interface{}) {
}

func (self MockLogger) Debug(message string) {
}

func (self MockLogger) Debugf(format string, values ...interface{}) {
}
