package logging

import "fmt"

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

// fmt.Stringify interface
func (self Level) String() string {
	switch self {
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
