package simple

import (
	"io"
	"strings"
	"time"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
)

const TIME_FORMAT = "2006/01/02 15:04:05.000"

type FormatFunc func(message string, id []string, level logging.Level, colorize bool) string

// FormatFunc signature
func DefaultFormat(message string, id []string, level logging.Level, colorize bool) string {
	var builder strings.Builder

	FormatTime(&builder)
	FormatLevel(&builder, level, true)
	builder.WriteRune(' ')
	FormatID(&builder, id)
	builder.WriteRune(' ')
	builder.WriteString(message)

	s := builder.String()

	if colorize {
		s = FormatColorize(s, level)
	}

	return s
}

func FormatTime(writer io.StringWriter) {
	writer.WriteString(time.Now().Format(TIME_FORMAT))
}

func FormatID(writer io.StringWriter, id []string) {
	writer.WriteString("[")
	length := len(id)
	switch length {
	case 0:
	case 1:
		writer.WriteString(id[0])
	default:
		last := length - 1
		for _, i := range id[:last] {
			writer.WriteString(i)
			writer.WriteString(".")
		}
		writer.WriteString(id[last])
	}
	writer.WriteString("]")
}

func FormatLevel(writer io.StringWriter, level logging.Level, align bool) {
	if align {
		switch level {
		case logging.Critical:
			writer.WriteString("  CRIT")
		case logging.Error:
			writer.WriteString(" ERROR")
		case logging.Warning:
			writer.WriteString("  WARN")
		case logging.Notice:
			writer.WriteString("  NOTE")
		case logging.Info:
			writer.WriteString("  INFO")
		case logging.Debug:
			writer.WriteString(" DEBUG")
		}
	} else {
		switch level {
		case logging.Critical:
			writer.WriteString("CRIT")
		case logging.Error:
			writer.WriteString("ERROR")
		case logging.Warning:
			writer.WriteString("WARN")
		case logging.Notice:
			writer.WriteString("NOTE")
		case logging.Info:
			writer.WriteString("INFO")
		case logging.Debug:
			writer.WriteString("DEBUG")
		}
	}
}

func FormatColorize(s string, level logging.Level) string {
	switch level {
	case logging.Critical:
		return terminal.ColorRed(s)
	case logging.Error:
		return terminal.ColorRed(s)
	case logging.Warning:
		return terminal.ColorYellow(s)
	case logging.Notice:
		return terminal.ColorCyan(s)
	case logging.Info:
		return terminal.ColorBlue(s)
	default:
		return s
	}
}
