package zerolog

import (
	"os"

	"github.com/rs/zerolog"
	logpkg "github.com/rs/zerolog/log"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
)

func init() {
	logging.SetBackend(NewBackend())
}

type Backend struct {
	loggers map[string]*Logger
}

func NewBackend() Backend {
	return Backend{
		loggers: make(map[string]*Logger),
	}
}

// logging.Backend interface

func (self Backend) Configure(verbosity int, path *string) {
	logpkg.Logger = zerolog.New(terminal.Stderr).With().Timestamp().Logger()
	for _, logger := range self.loggers {
		logger.Logger = logger.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func (self Backend) SetMaxLevel(name string, level logging.Level) {
	// TODO
}

func (self Backend) GetLogger(name string) logging.Logger {
	logger := Logger{logpkg.With().Str("name", name).Logger()}
	self.loggers[name] = &logger
	return &logger
}
