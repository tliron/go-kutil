package simple

import (
	"strings"
	"time"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
)

const TIME_FORMAT = "2006/01/02 15:04:05.000"

type FormatFunc func(message string, name string, level logging.Level, colorize bool) string

// FormatFunc signature
func DefaultFormat(message string, name string, level logging.Level, colorize bool) string {
	var builder strings.Builder

	builder.WriteString(time.Now().Format(TIME_FORMAT))

	switch level {
	case logging.Critical:
		builder.WriteString("  CRIT [")
	case logging.Error:
		builder.WriteString(" ERROR [")
	case logging.Warning:
		builder.WriteString("  WARN [")
	case logging.Notice:
		builder.WriteString("  NOTE [")
	case logging.Info:
		builder.WriteString("  INFO [")
	case logging.Debug:
		builder.WriteString(" DEBUG [")
	}

	builder.WriteString(name)
	builder.WriteString("] ")
	builder.WriteString(message)

	s := builder.String()

	if colorize {
		switch level {
		case logging.Critical:
			s = terminal.ColorRed(s)
		case logging.Error:
			s = terminal.ColorRed(s)
		case logging.Warning:
			s = terminal.ColorYellow(s)
		case logging.Notice:
			s = terminal.ColorCyan(s)
		case logging.Info:
			s = terminal.ColorBlue(s)
		}
	}

	return s
}
