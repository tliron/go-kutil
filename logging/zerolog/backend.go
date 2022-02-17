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
	Writer io.Writer

	hierarchy *logging.Hierarchy
}

func NewBackend() *Backend {
	return &Backend{
		hierarchy: logging.NewHierarchy(),
	}
}

// logging.Backend interface

func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := logging.VerbosityToMaxLevel(verbosity)

	if maxLevel == logging.None {
		self.Writer = io.Discard
		self.hierarchy.SetMaxLevel(nil, logging.None)
		logpkg.Logger = zerolog.New(self.Writer)
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				self.Writer = file
				logpkg.Logger = zerolog.New(self.Writer)
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.Writer = terminal.Stderr
			if terminal.Colorize {
				logpkg.Logger = zerolog.New(zerolog.ConsoleWriter{
					Out:        self.Writer,
					TimeFormat: TIME_FORMAT,
				})
			} else {
				logpkg.Logger = zerolog.New(zerolog.ConsoleWriter{
					Out:        self.Writer,
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

func (self *Backend) AllowLevel(id []string, level logging.Level) bool {
	return self.hierarchy.AllowLevel(id, level)
}

func (self *Backend) SetMaxLevel(id []string, level logging.Level) {
	self.hierarchy.SetMaxLevel(id, level)
}

func (self *Backend) NewMessage(id []string, level logging.Level, depth int) logging.Message {
	if self.AllowLevel(id, level) {
		logger := logpkg.With().Str("source", strings.Join(id, ".")).Logger()

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
