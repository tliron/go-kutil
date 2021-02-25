package klog

import (
	"k8s.io/klog/v2"
)

//
// Logger
//

type Logger struct {
	Prefix string
}

// logging.PartialLogger interface

func (self Logger) Critical(message string) {
	klog.ErrorDepth(3, self.Prefix+message)
}

func (self Logger) Error(message string) {
	klog.ErrorDepth(3, self.Prefix+message)
}

func (self Logger) Warning(message string) {
	klog.WarningDepth(3, self.Prefix+message)
}

func (self Logger) Notice(message string) {
	klog.InfoDepth(3, self.Prefix+message)
}

func (self Logger) Info(message string) {
	klog.InfoDepth(3, self.Prefix+message)
}

func (self Logger) Debug(message string) {
	klog.InfoDepth(3, self.Prefix+message)
}
