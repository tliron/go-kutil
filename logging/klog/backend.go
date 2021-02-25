package klog

import (
	"github.com/tliron/kutil/logging"
)

func init() {
	logging.SetBackend(NewBackend())
}

//
// Backend
//

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

// logging.Backend interface

func (self Backend) Configure(verbosity int, path *string) {
	// klog does its own configuration via klog.InitFlags
}

func (self Backend) SetMaxLevel(name string, level logging.Level) {
}

func (self Backend) GetLogger(name string) logging.Logger {
	return logging.FullLogger{
		Logger: Logger{
			Prefix: name + " ",
		},
	}
}
