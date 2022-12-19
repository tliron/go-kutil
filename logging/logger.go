package logging

import "fmt"

//
// Logger
//

type Logger interface {
	NewMessage(level Level, depth int) Message
	AllowLevel(level Level) bool
	SetMaxLevel(level Level)
	GetMaxLevel() Level

	// For unstructured logging

	Log(level Level, depth int, message string)
	Logf(level Level, depth int, format string, values ...any)

	Critical(message string)
	Criticalf(format string, values ...any)
	Error(message string)
	Errorf(format string, values ...any)
	Warning(message string)
	Warningf(format string, values ...any)
	Notice(message string)
	Noticef(format string, values ...any)
	Info(message string)
	Infof(format string, values ...any)
	Debug(message string)
	Debugf(format string, values ...any)
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

// Logger interface
func (self BackendLogger) AllowLevel(level Level) bool {
	return AllowLevel(self.name, level)
}

// Logger interface
func (self BackendLogger) SetMaxLevel(level Level) {
	SetMaxLevel(self.name, level)
}

// Logger interface
func (self BackendLogger) GetMaxLevel() Level {
	return GetMaxLevel(self.name)
}

// Logger interface
func (self BackendLogger) Log(level Level, depth int, message string) {
	if message_ := self.NewMessage(level, depth+1); message_ != nil {
		message_.Set("message", message)
		message_.Send()
	}
}

// Logger interface
func (self BackendLogger) Logf(level Level, depth int, format string, values ...any) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("message", fmt.Sprintf(format, values...))
		message.Send()
	}
}

// Logger interface
func (self BackendLogger) Critical(message string) {
	self.Log(Critical, 1, message)
}

// Logger interface
func (self BackendLogger) Criticalf(format string, values ...any) {
	self.Logf(Critical, 1, format, values...)
}

// Logger interface
func (self BackendLogger) Error(message string) {
	self.Log(Error, 1, message)
}

// Logger interface
func (self BackendLogger) Errorf(format string, values ...any) {
	self.Logf(Error, 1, format, values...)
}

// Logger interface
func (self BackendLogger) Warning(message string) {
	self.Log(Warning, 1, message)
}

// Logger interface
func (self BackendLogger) Warningf(format string, values ...any) {
	self.Logf(Warning, 1, format, values...)
}

// Logger interface
func (self BackendLogger) Notice(message string) {
	self.Log(Notice, 1, message)
}

// Logger interface
func (self BackendLogger) Noticef(format string, values ...any) {
	self.Logf(Notice, 1, format, values...)
}

// Logger interface
func (self BackendLogger) Info(message string) {
	self.Log(Info, 1, message)
}

// Logger interface
func (self BackendLogger) Infof(format string, values ...any) {
	self.Logf(Info, 1, format, values...)
}

// Logger interface
func (self BackendLogger) Debug(message string) {
	self.Log(Debug, 1, message)
}

// Logger interface
func (self BackendLogger) Debugf(format string, values ...any) {
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

// Logger interface
func (self ScopeLogger) AllowLevel(level Level) bool {
	return self.logger.AllowLevel(level)
}

// Logger interface
func (self ScopeLogger) SetMaxLevel(level Level) {
	self.logger.SetMaxLevel(level)
}

// Logger interface
func (self ScopeLogger) GetMaxLevel() Level {
	return self.logger.GetMaxLevel()
}

// Logger interface
func (self ScopeLogger) Log(level Level, depth int, message string) {
	if message_ := self.NewMessage(level, depth+1); message_ != nil {
		message_.Set("message", message)
		message_.Send()
	}
}

// Logger interface
func (self ScopeLogger) Logf(level Level, depth int, format string, values ...any) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("message", fmt.Sprintf(format, values...))
		message.Send()
	}
}

// Logger interface
func (self ScopeLogger) Critical(message string) {
	self.Log(Critical, 1, message)
}

// Logger interface
func (self ScopeLogger) Criticalf(format string, values ...any) {
	self.Logf(Critical, 1, format, values...)
}

// Logger interface
func (self ScopeLogger) Error(message string) {
	self.Log(Error, 1, message)
}

// Logger interface
func (self ScopeLogger) Errorf(format string, values ...any) {
	self.Logf(Error, 1, format, values...)
}

// Logger interface
func (self ScopeLogger) Warning(message string) {
	self.Log(Warning, 1, message)
}

// Logger interface
func (self ScopeLogger) Warningf(format string, values ...any) {
	self.Logf(Warning, 1, format, values...)
}

// Logger interface
func (self ScopeLogger) Notice(message string) {
	self.Log(Notice, 1, message)
}

// Logger interface
func (self ScopeLogger) Noticef(format string, values ...any) {
	self.Logf(Notice, 1, format, values...)
}

// Logger interface
func (self ScopeLogger) Info(message string) {
	self.Log(Info, 1, message)
}

// Logger interface
func (self ScopeLogger) Infof(format string, values ...any) {
	self.Logf(Info, 1, format, values...)
}

// Logger interface
func (self ScopeLogger) Debug(message string) {
	self.Log(Debug, 1, message)
}

// Logger interface
func (self ScopeLogger) Debugf(format string, values ...any) {
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

// Logger interface
func (self MockLogger) AllowLevel(level Level) bool {
	return false
}

// Logger interface
func (self MockLogger) SetMaxLevel(level Level) {
}

// Logger interface
func (self MockLogger) GetMaxLevel() Level {
	return None
}

// Logger interface
func (self MockLogger) Log(level Level, depth int, message string) {
}

// Logger interface
func (self MockLogger) Logf(level Level, depth int, format string, values ...any) {
}

// Logger interface
func (self MockLogger) Critical(message string) {
}

// Logger interface
func (self MockLogger) Criticalf(format string, values ...any) {
}

// Logger interface
func (self MockLogger) Error(message string) {
}

// Logger interface
func (self MockLogger) Errorf(format string, values ...any) {
}

// Logger interface
func (self MockLogger) Warning(message string) {
}

// Logger interface
func (self MockLogger) Warningf(format string, values ...any) {
}

// Logger interface
func (self MockLogger) Notice(message string) {
}

// Logger interface
func (self MockLogger) Noticef(format string, values ...any) {
}

// Logger interface
func (self MockLogger) Info(message string) {
}

// Logger interface
func (self MockLogger) Infof(format string, values ...any) {
}

// Logger interface
func (self MockLogger) Debug(message string) {
}

// Logger interface
func (self MockLogger) Debugf(format string, values ...any) {
}
