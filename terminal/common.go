package terminal

import (
	"io"
	"os"
)

var Stdout io.Writer = os.Stdout

var Stderr io.Writer = os.Stderr

var Quiet bool = false
