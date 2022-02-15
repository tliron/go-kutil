package format

import (
	"io"

	"github.com/fatih/color"
	yamllexer "github.com/goccy/go-yaml/lexer"
	yamlprinter "github.com/goccy/go-yaml/printer"
	"github.com/hokaccha/go-prettyjson"
	"github.com/kortschak/utter"
	"github.com/tliron/kutil/terminal"
)

func NewYAMLPrinter() yamlprinter.Printer {
	return yamlprinter.Printer{
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
}

func NewJSONFormatter() *prettyjson.Formatter {
	formatter := prettyjson.NewFormatter()
	formatter.Indent = terminal.IndentSpaces
	formatter.StringColor = color.New(color.FgBlue)
	formatter.NumberColor = color.New(color.FgMagenta)
	formatter.BoolColor = color.New(color.FgCyan)
	formatter.NullColor = color.New(color.FgCyan)
	formatter.KeyColor = color.New(color.FgGreen)
	return formatter
}

func PrettifyYAML(code string, writer io.Writer) error {
	tokens := yamllexer.Tokenize(code)
	printer := NewYAMLPrinter()
	if _, err := io.WriteString(writer, printer.PrintTokens(tokens)); err == nil {
		_, err := io.WriteString(writer, "\n")
		return err
	} else {
		return err
	}
}

func NewUtterConfig(indent string) *utter.ConfigState {
	var config = utter.NewDefaultConfig()
	config.Indent = indent
	config.SortKeys = true
	config.CommentPointers = true
	return config
}
