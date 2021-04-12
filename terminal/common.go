package terminal

import (
	"fmt"
	"io"
	"os"
)

var Stdout io.Writer = os.Stdout

var Stderr io.Writer = os.Stderr

var Stylize = NewStylist(false)

var Quiet bool = false

func Print(args ...interface{}) (int, error) {
	return fmt.Fprint(Stdout, args...)
}

func Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(Stdout, args...)
}

func Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(Stdout, format, args...)
}

func Eprint(args ...interface{}) (int, error) {
	return fmt.Fprint(Stderr, args...)
}

func Eprintln(args ...interface{}) (int, error) {
	return fmt.Fprintln(Stderr, args...)
}

func Eprintf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(Stderr, format, args...)
}
