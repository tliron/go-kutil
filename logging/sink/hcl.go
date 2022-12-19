package sink

import (
	"fmt"
	"io"
	logpkg "log"

	"github.com/hashicorp/go-hclog"
	"github.com/tliron/kutil/logging"
)

//
// Logger
//

type HCLogger struct {
	log  logging.Logger
	name string
	args []any
}

func NewHCLogger(name string, args []any) *HCLogger {
	return &HCLogger{
		log:  logging.GetLogger(name),
		name: name,
		args: args,
	}
}

// hclog.Logger interface
func (self *HCLogger) Log(level hclog.Level, msg string, args ...any) {
	self.sendMessage(level, msg, args)
}

// hclog.Logger interface
func (self *HCLogger) Trace(msg string, args ...any) {
	self.sendMessage(hclog.Trace, msg, args)
}

// hclog.Logger interface
func (self *HCLogger) Debug(msg string, args ...any) {
	self.sendMessage(hclog.Debug, msg, args)
}

// hclog.Logger interface
func (self *HCLogger) Info(msg string, args ...any) {
	self.sendMessage(hclog.Info, msg, args)
}

// hclog.Logger interface
func (self *HCLogger) Warn(msg string, args ...any) {
	self.sendMessage(hclog.Warn, msg, args)
}

// hclog.Logger interface
func (self *HCLogger) Error(msg string, args ...any) {
	self.sendMessage(hclog.Error, msg, args)
}

// hclog.Logger interface
func (self *HCLogger) IsTrace() bool {
	return self.log.AllowLevel(logging.Debug)
}

// hclog.Logger interface
func (self *HCLogger) IsDebug() bool {
	return self.log.AllowLevel(logging.Info)
}

// hclog.Logger interface
func (self *HCLogger) IsInfo() bool {
	return self.log.AllowLevel(logging.Notice)
}

// hclog.Logger interface
func (self *HCLogger) IsWarn() bool {
	return self.log.AllowLevel(logging.Warning)
}

// hclog.Logger interface
func (self *HCLogger) IsError() bool {
	return self.log.AllowLevel(logging.Error)
}

// hclog.Logger interface
func (self *HCLogger) ImpliedArgs() []any {
	return self.args
}

// hclog.Logger interface
func (self *HCLogger) With(args ...any) hclog.Logger {
	return NewHCLogger(self.name, args)
}

// hclog.Logger interface
func (self *HCLogger) Name() string {
	return self.name
}

// hclog.Logger interface
func (self *HCLogger) Named(name string) hclog.Logger {
	return NewHCLogger(self.name+"."+name, self.args)
}

// hclog.Logger interface
func (self *HCLogger) ResetNamed(name string) hclog.Logger {
	return NewHCLogger(name, self.args)
}

// hclog.Logger interface
func (self *HCLogger) SetLevel(level hclog.Level) {
	self.log.SetMaxLevel(toLevel(level))
}

// hclog.Logger interface
func (self *HCLogger) GetLevel() hclog.Level {
	return fromLevel(self.log.GetMaxLevel())
}

// hclog.Logger interface
func (self *HCLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *logpkg.Logger {
	// TODO
	return nil
}

// hclog.Logger interface
func (self *HCLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return logging.GetWriter()
}

// Utils

func (self *HCLogger) sendMessage(level hclog.Level, msg string, args []any) {
	if message := self.log.NewMessage(toLevel(level), 2); message != nil {
		message.Set("message", msg)
		args = append(self.args, args...)
		if length := len(args); length%2 == 0 {
			for i := 0; i < length; i += 2 {
				if key, ok := args[i].(string); ok {
					message.Set(key, args[i+1])
				}
			}
		}
		message.Send()
	}
}

func toLevel(level hclog.Level) logging.Level {
	switch level {
	case hclog.NoLevel:
		return logging.None
	case hclog.Trace:
		return logging.Debug
	case hclog.Debug:
		return logging.Info
	case hclog.Info:
		return logging.Notice
	case hclog.Warn:
		return logging.Warning
	case hclog.Error:
		return logging.Error
	default:
		panic(fmt.Sprintf("unsupported level: %d", level))
	}
}

func fromLevel(level logging.Level) hclog.Level {
	switch level {
	case logging.None:
		return hclog.NoLevel
	case logging.Critical:
		return hclog.Error
	case logging.Error:
		return hclog.Error
	case logging.Warning:
		return hclog.Warn
	case logging.Notice:
		return hclog.Info
	case logging.Info:
		return hclog.Debug
	case logging.Debug:
		return hclog.Trace
	default:
		panic(fmt.Sprintf("unsupported level: %d", level))
	}
}
