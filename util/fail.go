package util

import (
	"fmt"

	"github.com/tliron/kutil/terminal"
)

func Fail(message string) {
	if !terminal.Quiet {
		fmt.Fprintln(terminal.Stderr, terminal.StyleError(message))
	}
	Exit(1)
}

func Failf(f string, args ...interface{}) {
	Fail(fmt.Sprintf(f, args...))
}

func FailOnError(err error) {
	if err != nil {
		Fail(err.Error())
	}
}
