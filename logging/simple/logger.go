package simple

import (
	"io"
	"strings"

	"github.com/tliron/kutil/logging"
)

//
// Logger
//

type Logger struct {
	Backend *Backend
	Name    string
}

func (self Logger) Write(level logging.Level, message string) {
	if self.Backend.Allowed(self.Name, level) {
		if !strings.HasSuffix(message, "\n") {
			message += "\n"
		}

		message = self.Backend.format(message, self.Name, level, self.Backend.colorize)
		io.WriteString(self.Backend.writer, message)
	}
}

// logging.PartialLogger interface

func (self Logger) Critical(message string) {
	self.Write(logging.Critical, message)
}

func (self Logger) Error(message string) {
	self.Write(logging.Error, message)
}

func (self Logger) Warning(message string) {
	self.Write(logging.Warning, message)
}

func (self Logger) Notice(message string) {
	self.Write(logging.Notice, message)
}

func (self Logger) Info(message string) {
	self.Write(logging.Info, message)
}

func (self Logger) Debug(message string) {
	self.Write(logging.Debug, message)
}
