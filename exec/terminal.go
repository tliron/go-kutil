package exec

import (
	"os"

	"golang.org/x/term"
)

//
// Terminal
//

type Terminal struct {
	InitialSize *Size
	Resize      chan Size // receive from this

	sigwinch  chan os.Signal
	termState *term.State
}
