package logging

import (
	"fmt"
	"strings"
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

func AllowLevel(id []string, level Level) bool {
	if backend != nil {
		return backend.AllowLevel(id, level)
	} else {
		return false
	}
}

func SetMaxLevel(id []string, level Level) {
	if backend != nil {
		backend.SetMaxLevel(id, level)
	}
}

func NewMessage(id []string, level Level, depth int) Message {
	if backend != nil {
		return backend.NewMessage(id, level, depth)
	} else {
		return nil
	}
}

// Unstructured wrappers

func GetLogger(name string) Logger {
	return NewBackendLogger(strings.Split(name, "."))
}

func GetLoggerf(format string, values ...interface{}) Logger {
	return GetLogger(fmt.Sprintf(format, values...))
}
