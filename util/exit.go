package util

import (
	"os"
	"sync"
)

// Inspired by: https://github.com/tebeka/atexit

var exitEntries []exitEntry
var exitNextHandle ExitFunctionHandle
var exitLock sync.RWMutex

type exitEntry struct {
	function func()
	handle   ExitFunctionHandle
}

func OnExit(exitFunction func()) ExitFunctionHandle {
	exitLock.Lock()
	defer exitLock.Unlock()

	handle := exitNextHandle
	exitNextHandle++

	exitEntries = append(exitEntries, exitEntry{
		function: exitFunction,
		handle:   handle,
	})

	return handle
}

func Exit(code int) {
	exitLock.RLock()
	defer exitLock.RUnlock()

	for _, exitEntry := range exitEntries {
		exitEntry.function()
	}

	os.Exit(code)
}

//
// ExitFunctionHandle
//

type ExitFunctionHandle uint64

func (self ExitFunctionHandle) Cancel() {
	exitLock.Lock()
	defer exitLock.Unlock()

	for index, exitEntry := range exitEntries {
		if exitEntry.handle == self {
			exitEntries = append(exitEntries[:index], exitEntries[index+1:]...)
			break
		}
	}
}
