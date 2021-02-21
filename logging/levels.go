package logging

import (
	loggingpkg "github.com/op/go-logging"
)

//
// Level
//

type Level int

const (
	Critical Level = 1
	Error    Level = 2
	Warning  Level = 3
	Notice   Level = 4
	Info     Level = 5
	Debug    Level = 6
)

func SetMaxLevel(name string, level Level) {
	loggingpkg.SetLevel(loggingpkg.Level(level-1), name)
}
