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
	logging.SetBackend(NewBackend())
}

const LOG_FILE_WRITE_PERMISSIONS = 0600

const TIME_FORMAT = "2006/01/02 15:04:05.000"

//
// Backend
//

type Backend struct {
	writer  io.Writer
	loggers map[string]*Logger
}

func NewBackend() Backend {
	return Backend{
		loggers: make(map[string]*Logger),
	}
}

// logging.Backend interface

func (self Backend) Configure(verbosity int, path *string) {
	if verbosity == -1 {
		self.writer = io.Discard
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
			self.writer = terminal.Stderr
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

		switch verbosity {
		case 0:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case 1:
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
		logpkg.Logger = logpkg.With().Timestamp().Logger()
	}
}

func (self Backend) GetWriter() io.Writer {
	return self.writer
}

func (self Backend) SetMaxLevel(name string, level logging.Level) {
	if strings.HasSuffix(name, "*") {
		if len(name) == 1 {
			zerolog.SetGlobalLevel(toLevel(level))
		} else {
			prefix := name[:len(name)-1]
			for name_, logger := range self.loggers {
				if strings.HasPrefix(name_, prefix) {
					logger.SetMaxLevel(level)
				}
			}
		}
	} else {
		if logger, ok := self.loggers[name]; ok {
			logger.SetMaxLevel(level)
		}
	}
}

func (self Backend) GetLogger(name string) logging.Logger {
	// TODO: new loggers won't respect SetMaxLevel!
	logger := &Logger{
		Logger: logpkg.With().Str("logger", name).Logger(),
	}
	self.loggers[name] = logger
	return logger
}

// Utils

func toLevel(level logging.Level) zerolog.Level {
	switch level {
	case logging.Critical:
		return zerolog.ErrorLevel
	case logging.Error:
		return zerolog.ErrorLevel
	case logging.Warning:
		return zerolog.WarnLevel
	case logging.Notice:
		return zerolog.InfoLevel
	case logging.Info:
		return zerolog.DebugLevel
	case logging.Debug:
		return zerolog.TraceLevel
	default:
		panic(fmt.Sprintf("unsupported level: %d", level))
	}
}
