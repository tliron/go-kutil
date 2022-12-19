package journal

import (
	"fmt"
	"io"
	"strings"

	"github.com/coreos/go-systemd/journal"
	"github.com/tliron/kutil/logging"
)

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	logging.SetBackend(backend)
}

//
// Backend
//

type Backend struct {
	hierarchy *logging.Hierarchy
	writer    io.Writer
}

func NewBackend() *Backend {
	return &Backend{
		hierarchy: logging.NewMaxLevelHierarchy(),
	}
}

// logging.Backend interface
func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := logging.VerbosityToMaxLevel(verbosity)

	if maxLevel == logging.None {
		self.writer = io.Discard
		self.hierarchy.SetMaxLevel(nil, logging.None)
	} else {
		self.writer = JournalWriter{}
		self.hierarchy.SetMaxLevel(nil, maxLevel)
	}
}

// logging.Backend interface
func (self *Backend) GetWriter() io.Writer {
	return self.writer
}

// logging.Backend interface
func (self *Backend) NewMessage(name []string, level logging.Level, depth int) logging.Message {
	if self.AllowLevel(name, level) {
		var priority journal.Priority
		switch level {
		case logging.Critical:
			priority = journal.PriCrit
		case logging.Error:
			priority = journal.PriErr
		case logging.Warning:
			priority = journal.PriWarning
		case logging.Notice:
			priority = journal.PriNotice
		case logging.Info:
			priority = journal.PriInfo
		case logging.Debug:
			priority = journal.PriDebug
		default:
			panic(fmt.Sprintf("unsupported level: %d", level))
		}

		var prefix string
		if name := strings.Join(name, "."); len(name) > 0 {
			prefix = "[" + name + "] "
		}

		return NewMessage(priority, prefix)
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

// logging.Backend interface
func (self *Backend) GetMaxLevel(name []string) logging.Level {
	return self.hierarchy.GetMaxLevel(name)
}
