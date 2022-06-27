package simple

import (
	"io"
	"os"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

const LOG_FILE_WRITE_PERMISSIONS = 0600

const BUFFER_SIZE = 10000

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	logging.SetBackend(backend)
}

//
// Backend
//

type Backend struct {
	Writer   io.Writer
	Format   FormatFunc
	Buffered bool

	colorize  bool
	hierarchy *logging.Hierarchy
}

func NewBackend() *Backend {
	return &Backend{
		Format:    DefaultFormat,
		Buffered:  true,
		hierarchy: logging.NewMaxLevelHierarchy(),
	}
}

// logging.Backend interface
func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := logging.VerbosityToMaxLevel(verbosity)

	if maxLevel == logging.None {
		self.Writer = io.Discard
		self.hierarchy.SetMaxLevel(nil, logging.None)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				if self.Buffered {
					writer := util.NewBufferedWriter(file, BUFFER_SIZE)
					util.OnExitError(writer.Close)
					self.Writer = writer
				} else {
					self.Writer = util.NewSyncedWriter(file)
				}
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.colorize = terminal.Colorize
			if self.Buffered {
				writer := util.NewBufferedWriter(terminal.Stderr, BUFFER_SIZE)
				util.OnExitError(writer.Close)
				self.Writer = writer
			} else {
				self.Writer = util.NewSyncedWriter(terminal.Stderr)
			}
		}

		self.hierarchy.SetMaxLevel(nil, maxLevel)
	}
}

// logging.Backend interface
func (self *Backend) GetWriter() io.Writer {
	return self.Writer
}

// logging.Backend interface
func (self *Backend) NewMessage(name []string, level logging.Level, depth int) logging.Message {
	if self.AllowLevel(name, level) {
		return logging.NewUnstructuredMessage(func(message string) {
			message = self.Format(message, name, level, self.colorize)
			io.WriteString(self.Writer, message+"\n")
		})
	} else {
		return nil
	}
}

// logging.Backend interface
func (self *Backend) AllowLevel(name []string, level logging.Level) bool {
	return self.hierarchy.AllowLevel(name, level)
}

// logging.Backend interface
func (self *Backend) SetMaxLevel(name []string, level logging.Level) {
	self.hierarchy.SetMaxLevel(name, level)
}
