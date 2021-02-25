package logging

import (
	"fmt"
)

var backend Backend

func SetBackend(backend_ Backend) {
	backend = backend_
}

func Configure(verbosity int, path *string) {
	backend.Configure(verbosity, path)
}

func SetMaxLevel(name string, level Level) {
	backend.SetMaxLevel(name, level)
}

func GetLogger(name string) Logger {
	return &LazyLogger{
		Name: name,
	}
}

func GetLoggerf(format string, arguments ...interface{}) Logger {
	return GetLogger(fmt.Sprintf(format, arguments...))
}
