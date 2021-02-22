package simple

import (
	"fmt"
	"io"
	"strings"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
)

//
// Logger
//

type Logger struct {
	Backend *Backend
	Name    string
}

func (self *Logger) Write(level logging.Level, message string) {
	if self.Backend.Allowed(self.Name, level) {
		if !strings.HasSuffix(message, "\n") {
			message += "\n"
		}

		message = Format(message, self.Name, level, terminal.Colorize)

		io.WriteString(self.Backend.writer, message)
	}
}

func (self *Logger) Writef(level logging.Level, format string, values ...interface{}) {
	self.Write(level, fmt.Sprintf(format, values...))
}

// logging.Logger interface

func (self *Logger) Critical(message string) {
	self.Write(logging.Critical, message)
}

func (self *Logger) Criticalf(format string, values ...interface{}) {
	self.Writef(logging.Critical, format, values...)
}

func (self *Logger) Error(message string) {
	self.Write(logging.Error, message)
}

func (self *Logger) Errorf(format string, values ...interface{}) {
	self.Writef(logging.Error, format, values...)
}

func (self *Logger) Warning(message string) {
	self.Write(logging.Warning, message)
}

func (self *Logger) Warningf(format string, values ...interface{}) {
	self.Writef(logging.Warning, format, values...)
}

func (self *Logger) Notice(message string) {
	self.Write(logging.Notice, message)
}

func (self *Logger) Noticef(format string, values ...interface{}) {
	self.Writef(logging.Notice, format, values...)
}

func (self *Logger) Info(message string) {
	self.Write(logging.Info, message)
}

func (self *Logger) Infof(format string, values ...interface{}) {
	self.Writef(logging.Info, format, values...)
}

func (self *Logger) Debug(message string) {
	self.Write(logging.Debug, message)
}

func (self *Logger) Debugf(format string, values ...interface{}) {
	self.Writef(logging.Debug, format, values...)
}
