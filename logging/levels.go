package logging

import (
	"fmt"
)

//
// Level
//

type Level int

const (
	None     Level = 0
	Critical Level = 1
	Error    Level = 2
	Warning  Level = 3
	Notice   Level = 4
	Info     Level = 5
	Debug    Level = 6
)

// fmt.Stringify interface
func (self Level) String() string {
	switch self {
	case None:
		return "None"
	case Critical:
		return "Critical"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Notice:
		return "Notice"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	default:
		panic(fmt.Sprintf("unsupported level: %d", self))
	}
}

func VerbosityToMaxLevel(verbosity int) Level {
	if verbosity < 0 {
		return None
	} else {
		switch verbosity {
		case 0:
			return Notice
		case 1:
			return Info
		default:
			return Debug
		}
	}
}
