package transcribe

import (
	"io"

	"github.com/fatih/color"
	yamllexer "github.com/goccy/go-yaml/lexer"
	yamlprinter "github.com/goccy/go-yaml/printer"
	"github.com/hokaccha/go-prettyjson"
	"github.com/tliron/kutil/terminal"
)

var yamlColorPrinter = yamlprinter.Printer{
	String: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.BlueCode,
			Suffix: terminal.ResetCode,
		}
	},
	Number: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.MagentaCode,
			Suffix: terminal.ResetCode,
		}
	},
	Bool: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.CyanCode,
			Suffix: terminal.ResetCode,
		}
	},

	MapKey: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.GreenCode,
			Suffix: terminal.ResetCode,
		}
	},

	Anchor: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.RedCode,
			Suffix: terminal.ResetCode,
		}
	},
	Alias: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.YellowCode,
			Suffix: terminal.ResetCode,
		}
	},
}

func ColorizeYAML(code string, writer io.Writer) error {
	tokens := yamllexer.Tokenize(code)
	if _, err := io.WriteString(writer, yamlColorPrinter.PrintTokens(tokens)); err == nil {
		_, err := io.WriteString(writer, "\n")
		return err
	} else {
		return err
	}
}

func NewJSONColorFormatter(indent int) *prettyjson.Formatter {
	return &prettyjson.Formatter{
		KeyColor:        color.New(color.FgGreen),
		StringColor:     color.New(color.FgBlue),
		BoolColor:       color.New(color.FgCyan),
		NumberColor:     color.New(color.FgMagenta),
		NullColor:       color.New(color.FgCyan),
		StringMaxLength: 0,
		DisabledColor:   false,
		Indent:          indent,
		Newline:         "\n",
	}
}
