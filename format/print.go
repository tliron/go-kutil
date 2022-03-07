package format

import (
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func Print(value any, format string, writer io.Writer, strict bool, pretty bool) error {
	// Special handling for strings (ignore format)
	if s, ok := value.(string); ok {
		if _, err := io.WriteString(writer, s); err == nil {
			if pretty {
				return util.WriteNewline(writer)
			} else {
				return nil
			}
		} else {
			return err
		}
	}

	// Special handling for etree (ignore format)
	if xmlDocument, ok := value.(*etree.Document); ok {
		return PrintXMLDocument(xmlDocument, writer, pretty)
	}

	switch format {
	case "yaml", "":
		return PrintYAML(value, writer, strict, pretty)

	case "json":
		return PrintJSON(value, writer, pretty)

	case "cjson":
		return PrintCompatibleJSON(value, writer, pretty)

	case "xml":
		return PrintCompatibleXML(value, writer, pretty)

	case "cbor":
		return PrintCBOR(value, writer)

	case "go":
		return PrintGo(value, writer, pretty)

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintYAML(value any, writer io.Writer, strict bool, pretty bool) error {
	if pretty && terminal.Colorize {
		if code, err := EncodeYAML(value, terminal.Indent, strict); err == nil {
			return PrettifyYAML(code, writer)
		} else {
			return err
		}
	} else {
		return WriteYAML(value, writer, "  ", strict)
	}
}

func PrintJSON(value any, writer io.Writer, pretty bool) error {
	if pretty {
		if terminal.Colorize {
			formatter := NewJSONFormatter()
			if bytes, err := formatter.Marshal(value); err == nil {
				if _, err := writer.Write(bytes); err == nil {
					return util.WriteNewline(writer)
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			return WriteJSON(value, writer, terminal.Indent)
		}
	} else {
		return WriteJSON(value, writer, "")
	}
}

func PrintCompatibleJSON(value any, writer io.Writer, pretty bool) error {
	return PrintJSON(ard.EnsureCompatibleJSON(value), writer, pretty)
}

func PrintCompatibleXML(value any, writer io.Writer, pretty bool) error {
	indent := ""
	if pretty {
		indent = terminal.Indent
	}
	if err := WriteCompatibleXML(value, writer, indent); err == nil {
		if pretty {
			return util.WriteNewline(writer)
		} else {
			return nil
		}
	} else {
		return err
	}
}

func PrintXMLDocument(xmlDocument *etree.Document, writer io.Writer, pretty bool) error {
	if pretty {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		xmlDocument.Indent(0)
	}
	_, err := xmlDocument.WriteTo(writer)
	return err
}

func PrintCBOR(value any, writer io.Writer) error {
	return WriteCBOR(value, writer)
}

func PrintGo(value any, writer io.Writer, pretty bool) error {
	if pretty {
		return WriteGo(value, writer, terminal.Indent)
	} else {
		return WriteGo(value, writer, "")
	}
}
