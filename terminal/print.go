package terminal

import (
	"fmt"
)

// stdout

func Print(args ...interface{}) (int, error) {
	return fmt.Fprint(Stdout, args...)
}

func Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(Stdout, args...)
}

func Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(Stdout, format, args...)
}

// stderr

func Eprint(args ...interface{}) (int, error) {
	return fmt.Fprint(Stderr, args...)
}

func Eprintln(args ...interface{}) (int, error) {
	return fmt.Fprintln(Stderr, args...)
}

func Eprintf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(Stderr, format, args...)
}
