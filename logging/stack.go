package logging

import (
	"fmt"
	"runtime"
	"strings"

	loggingpkg "github.com/op/go-logging"
)

func LogStack(log *loggingpkg.Logger, message string, skip int) {
	// See: https://golang.org/src/runtime/example_test.go

	var builder strings.Builder
	fmt.Fprintf(&builder, "%s\n", message)

	pc := make([]uintptr, 30)
	if n := runtime.Callers(skip+2, pc); n > 0 {
		pc = pc[:n]
		frames := runtime.CallersFrames(pc)
		for {
			frame, more := frames.Next()
			fn := runtime.FuncForPC(frame.PC)
			fmt.Fprintf(&builder, "%s\n  %s:%d\n", fn.Name(), frame.File, frame.Line)

			if !more {
				break
			}
		}
	}

	log.Critical(strings.TrimRight(builder.String(), "\n"))
}
