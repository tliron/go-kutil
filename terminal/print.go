package terminal

import (
	"fmt"
)

// stdout

func Print(args ...any) (int, error) {
	return fmt.Fprint(Stdout, args...)
}

func Println(args ...any) (int, error) {
	return fmt.Fprintln(Stdout, args...)
}

func Printf(format string, args ...any) (int, error) {
	return fmt.Fprintf(Stdout, format, args...)
}

// stderr

func Eprint(args ...any) (int, error) {
	return fmt.Fprint(Stderr, args...)
}

func Eprintln(args ...any) (int, error) {
	return fmt.Fprintln(Stderr, args...)
}

func Eprintf(format string, args ...any) (int, error) {
	return fmt.Fprintf(Stderr, format, args...)
}
