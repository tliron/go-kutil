package logging

import (
	"fmt"
	"io"
)

var backend Backend

func SetBackend(backend_ Backend) {
	backend = backend_
}

func Configure(verbosity int, path *string) {
	if backend != nil {
		backend.Configure(verbosity, path)
	}
}

func GetWriter() io.Writer {
	if backend != nil {
		return backend.GetWriter()
	} else {
		return io.Discard
	}
}

func SetMaxLevel(name string, level Level) {
	if backend != nil {
		backend.SetMaxLevel(name, level)
	}
}

func GetLogger(name string) Logger {
	return &LazyLogger{
		Name: name,
	}
}

func GetLoggerf(format string, arguments ...interface{}) Logger {
	return GetLogger(fmt.Sprintf(format, arguments...))
}
