package simple

import (
	"io"
	"os"
	"sort"
	"strings"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	logging.SetBackend(backend)
}

const LOG_FILE_WRITE_PERMISSIONS = 0600

const BUFFER_SIZE = 10000

//
// Backend
//

type Backend struct {
	writer          io.Writer
	maxLevel        logging.Level
	maxLevels       map[string]logging.Level
	prefixMaxLevels []prefixLevel
	format          FormatFunc
	colorize        bool
	buffered        bool
}

type prefixLevel struct {
	prefix string
	level  logging.Level
}

func NewBackend() *Backend {
	return &Backend{
		maxLevels: make(map[string]logging.Level),
		format:    DefaultFormat,
		buffered:  true,
	}
}

func (self *Backend) SetFormat(format FormatFunc) {
	self.format = format
}

func (self *Backend) SetBuffered(buffered bool) {
	self.buffered = buffered
}

func (self *Backend) GetMaxLevel(name string) logging.Level {
	for _, prefixLevel := range self.prefixMaxLevels {
		if strings.HasPrefix(name, prefixLevel.prefix) {
			return prefixLevel.level
		}
	}

	if level, ok := self.maxLevels[name]; ok {
		return level
	}

	return self.maxLevel
}

func (self *Backend) Allowed(name string, level logging.Level) bool {
	return level <= self.GetMaxLevel(name)
}

// logging.Backend interface

func (self *Backend) Configure(verbosity int, path *string) {
	if verbosity == -1 {
		self.writer = io.Discard
		self.maxLevel = logging.Level(0)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				if self.buffered {
					writer := util.NewBufferedWriter(file, BUFFER_SIZE)
					writer.CloseOnExit()
					self.writer = writer
				} else {
					self.writer = util.NewSyncedWriter(file)
				}
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.colorize = terminal.Colorize
			if self.buffered {
				writer := util.NewBufferedWriter(terminal.Stderr, BUFFER_SIZE)
				writer.CloseOnExit()
				self.writer = writer
			} else {
				self.writer = util.NewSyncedWriter(terminal.Stderr)
			}
		}

		verbosity += 4 // our 0 verbosity is max level NOTICE (4)
		if verbosity > 6 {
			verbosity = 6
		}

		self.maxLevel = logging.Level(verbosity)
	}
}

func (self *Backend) GetWriter() io.Writer {
	return self.writer
}

func (self *Backend) SetMaxLevel(name string, level logging.Level) {
	if strings.HasSuffix(name, "*") {
		if len(name) == 1 {
			self.maxLevel = level
		} else {
			self.prefixMaxLevels = append(self.prefixMaxLevels, prefixLevel{
				prefix: name[:len(name)-1],
				level:  level,
			})

			// Sort in reverse so that the more specific (=longer) prefixes come first
			sort.Slice(self.prefixMaxLevels, func(i int, j int) bool {
				return strings.Compare(self.prefixMaxLevels[i].prefix, self.prefixMaxLevels[j].prefix) == 1
			})
		}
	} else {
		self.maxLevels[name] = level
	}
}

func (self *Backend) GetLogger(name string) logging.Logger {
	return logging.FullLogger{
		Logger: Logger{
			Backend: self,
			Name:    name,
		},
	}
}
