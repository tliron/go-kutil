package gologging

import (
	"io"
	"os"

	loggingpkg "github.com/op/go-logging"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	logging.SetBackend(NewBackend())
}

const LOG_FILE_WRITE_PERMISSIONS = 0600

var plainFormatter = loggingpkg.MustStringFormatter(
	`%{time:2006/01/02 15:04:05.000} %{level:8.8s} [%{module}] %{message}`,
)

var colorFormatter = loggingpkg.MustStringFormatter(
	`%{color}%{time:2006/01/02 15:04:05.000} %{level:8.8s} [%{module}] %{message}%{color:reset}`,
)

//
// Backend
//

type Backend struct {
	writer io.Writer
}

func NewBackend() Backend {
	return Backend{}
}

// logging.Backend interface

func (self Backend) Configure(verbosity int, path *string) {
	var backend *loggingpkg.LogBackend

	if verbosity == -1 {
		self.writer = io.Discard
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExit(func() {
					file.Close()
				})
				self.writer = file
				loggingpkg.SetFormatter(plainFormatter)
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.writer = terminal.Stderr
			if terminal.Colorize {
				loggingpkg.SetFormatter(colorFormatter)
			} else {
				loggingpkg.SetFormatter(plainFormatter)
			}
		}

		verbosity += 3 // our 0 verbosity is max level NOTICE (3)
		if verbosity > 5 {
			verbosity = 5
		}
	}

	backend = loggingpkg.NewLogBackend(self.writer, "", 0)
	leveledBackend := loggingpkg.AddModuleLevel(backend)
	leveledBackend.SetLevel(loggingpkg.Level(verbosity), "")
	loggingpkg.SetBackend(NewPrefixLeveledBackend(leveledBackend))
}

func (self Backend) GetWriter() io.Writer {
	return self.writer
}

func (self Backend) SetMaxLevel(name string, level logging.Level) {
	loggingpkg.SetLevel(loggingpkg.Level(level-1), name)
}

func (self Backend) GetLogger(name string) logging.Logger {
	return Logger{
		Logger: loggingpkg.MustGetLogger(name),
	}
}
