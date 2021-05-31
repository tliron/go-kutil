package util

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	length := len(exitEntries)
	for index := range exitEntries {
		exitEntry := exitEntries[length-index-1]
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "panic during exit: %s\n", r)
				}
			}()

			exitEntry.function()
		}()
	}

	exitLock.RUnlock()

	os.Exit(code)
}

func ExitOnSIGTERM() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Exit(1)
	}()
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
