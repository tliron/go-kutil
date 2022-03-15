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

func (self *HCLogger) Trace(msg string, args ...any) {
	self.sendMessage(hclog.Trace, msg, args)
}

func (self *HCLogger) Debug(msg string, args ...any) {
	self.sendMessage(hclog.Debug, msg, args)
}

func (self *HCLogger) Info(msg string, args ...any) {
	self.sendMessage(hclog.Info, msg, args)
}

func (self *HCLogger) Warn(msg string, args ...any) {
	self.sendMessage(hclog.Warn, msg, args)
}

func (self *HCLogger) Error(msg string, args ...any) {
	self.sendMessage(hclog.Error, msg, args)
}

func (self *HCLogger) IsTrace() bool {
	return self.log.AllowLevel(logging.Debug)
}

func (self *HCLogger) IsDebug() bool {
	return self.log.AllowLevel(logging.Info)
}

func (self *HCLogger) IsInfo() bool {
	return self.log.AllowLevel(logging.Notice)
}

func (self *HCLogger) IsWarn() bool {
	return self.log.AllowLevel(logging.Warning)
}

func (self *HCLogger) IsError() bool {
	return self.log.AllowLevel(logging.Error)
}

func (self *HCLogger) ImpliedArgs() []any {
	return self.args
}

func (self *HCLogger) With(args ...any) hclog.Logger {
	return NewHCLogger(self.name, args)
}

func (self *HCLogger) Name() string {
	return self.name
}

func (self *HCLogger) Named(name string) hclog.Logger {
	return NewHCLogger(self.name+"."+name, self.args)
}

func (self *HCLogger) ResetNamed(name string) hclog.Logger {
	return NewHCLogger(name, self.args)
}

func (self *HCLogger) SetLevel(level hclog.Level) {
	self.log.SetMaxLevel(toLevel(level))
}

func (self *HCLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *logpkg.Logger {
	// TODO
	return nil
}

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
