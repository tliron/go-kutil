package simple

import (
	"strings"
	"time"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
)

const TIME_FORMAT = "2006/01/02 15:04:05.000"

func Format(message string, name string, level logging.Level, colorize bool) string {
	var builder strings.Builder

	builder.WriteString(time.Now().Format(TIME_FORMAT))

	switch level {
	case logging.Critical:
		builder.WriteString(" CRITICAL [")
	case logging.Error:
		builder.WriteString("    ERROR [")
	case logging.Warning:
		builder.WriteString("  WARNING [")
	case logging.Notice:
		builder.WriteString("   NOTICE [")
	case logging.Info:
		builder.WriteString("     INFO [")
	case logging.Debug:
		builder.WriteString("    DEBUG [")
	}

	builder.WriteString(name)
	builder.WriteString("] ")
	builder.WriteString(message)

	s := builder.String()

	if colorize {
		switch level {
		case logging.Critical:
			s = terminal.ColorError(s)
		case logging.Error:
			s = terminal.ColorError(s)
		case logging.Warning:
			s = terminal.ColorValue(s)
		case logging.Notice:
			s = terminal.ColorName(s)
		case logging.Info:
			s = terminal.ColorName(s)
		}
	}

	return s
}
