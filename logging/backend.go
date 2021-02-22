package logging

import "fmt"

type Backend interface {
	// If "path" is nil will log to stdout, colorized if possible
	// The default "verbosity" 0 will log criticals, errors, warnings, and notices.
	// "verbosity" 1 will add infos. "verbosity" 2 will add debugs.
	// Set "verbostiy" to -1 to disable the log.
	Configure(verbosity int, path *string)

	SetMaxLevel(name string, level Level)
	GetLogger(name string) Logger
}

var current Backend

func SetBackend(backend Backend) {
	if current == nil {
		current = backend
	}
}

func Configure(verbosity int, path *string) {
	current.Configure(verbosity, path)
}

func SetMaxLevel(name string, level Level) {
	current.SetMaxLevel(name, level)
}

func GetLogger(name string) Logger {
	return current.GetLogger(name)
}

func GetLoggerf(format string, arguments ...interface{}) Logger {
	return current.GetLogger(fmt.Sprintf(format, arguments...))
}
