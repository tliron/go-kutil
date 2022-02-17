package logging

import "fmt"

//
// Backend
//

type Backend interface {
	// If "path" is nil will log to stdout, colorized if possible
	// The default "verbosity" 0 will log criticals, errors, warnings, and notices.
	// "verbosity" 1 will add infos. "verbosity" 2 will add debugs.
	// Set "verbostiy" to -1 to disable the log.
	Configure(verbosity int, path *string)

	AllowLevel(id []string, level Level) bool
	SetMaxLevel(id []string, level Level)

	// May return nil
	NewMessage(id []string, level Level, depth int) Message
}

//
// BackendLogger
//

type BackendLogger struct {
	id []string
}

func NewBackendLogger(id []string) BackendLogger {
	return BackendLogger{id: id}
}

// Logger interface

func (self BackendLogger) NewMessage(level Level, depth int) Message {
	return NewMessage(self.id, level, depth)
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

func (self BackendLogger) AllowLevel(level Level) bool {
	return AllowLevel(self.id, level)
}

func (self BackendLogger) SetMaxLevel(level Level) {
	SetMaxLevel(self.id, level)
}
