package zerolog

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	logpkg "github.com/rs/zerolog/log"
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

const TIME_FORMAT = "2006/01/02 15:04:05.000"

//
// Backend
//

// Note: using kutil to wrap zerolog will circumvent its primary optimization, which is for high
// performance and low resource use due to aggressively avoiding allocations. If you need that
// optimization then you should use zerolog's API directly.

type Backend struct {
	writer    io.Writer
	hierarchy *logging.Hierarchy
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
		logpkg.Logger = zerolog.New(self.writer)
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				self.writer = file
				logpkg.Logger = zerolog.New(self.writer)
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.writer = os.Stderr
			if terminal.Colorize {
				logpkg.Logger = zerolog.New(zerolog.ConsoleWriter{
					Out:        self.writer,
					TimeFormat: TIME_FORMAT,
				})
			} else {
				logpkg.Logger = zerolog.New(zerolog.ConsoleWriter{
					Out:        self.writer,
					TimeFormat: TIME_FORMAT,
					NoColor:    true,
				})
			}
		}

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
		logpkg.Logger = logpkg.With().Timestamp().Logger()

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
		context := logpkg.With()
		if name := strings.Join(name, "."); len(name) > 0 {
			context = context.Str("name", name)
		}
		logger := context.Logger()

		var event *zerolog.Event
		switch level {
		case logging.Critical:
			event = logger.Error()
		case logging.Error:
			event = logger.Error()
		case logging.Warning:
			event = logger.Warn()
		case logging.Notice:
			event = logger.Info()
		case logging.Info:
			event = logger.Debug()
		case logging.Debug:
			event = logger.Trace()
		default:
			panic(fmt.Sprintf("unsupported level: %d", level))
		}

		return NewMessage(event)
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
