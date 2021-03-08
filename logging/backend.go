package logging

import (
	"io"
)

type Backend interface {
	// If "path" is nil will log to stdout, colorized if possible
	// The default "verbosity" 0 will log criticals, errors, warnings, and notices.
	// "verbosity" 1 will add infos. "verbosity" 2 will add debugs.
	// Set "verbostiy" to -1 to disable the log.
	Configure(verbosity int, path *string)

	GetWriter() io.Writer
	SetMaxLevel(name string, level Level)
	GetLogger(name string) Logger
}
