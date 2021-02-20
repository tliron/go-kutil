package format

import (
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/hokaccha/go-prettyjson"
	"github.com/tliron/kutil/terminal"
)

func Print(value interface{}, format string, writer io.Writer, strict bool, pretty bool) error {
	// Special handling for strings (ignore format)
	if s, ok := value.(string); ok {
		if pretty {
			s += "\n"
		}
		_, err := fmt.Fprint(writer, s)
		return err
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

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintYAML(value interface{}, writer io.Writer, strict bool, pretty bool) error {
	indent := "  "
	if pretty {
		indent = terminal.Indent
	}
	return WriteYAML(value, writer, indent, strict)
}

func PrintJSON(value interface{}, writer io.Writer, pretty bool) error {
	if pretty {
		prettyJsonFormatter := prettyjson.NewFormatter()
		prettyJsonFormatter.Indent = terminal.IndentSpaces
		bytes, err := prettyJsonFormatter.Marshal(value)
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, "%s\n", bytes)
	} else {
		return WriteJSON(value, writer, "")
	}
	return nil
}

func PrintCompatibleJSON(value interface{}, writer io.Writer, pretty bool) error {
	return PrintJSON(ToCompatibleJSON(value), writer, pretty)
}

func PrintCompatibleXML(value interface{}, writer io.Writer, pretty bool) error {
	indent := ""
	if pretty {
		indent = terminal.Indent
	}
	if err := WriteCompatibleXML(value, writer, indent); err != nil {
		return err
	}
	if pretty {
		fmt.Fprintln(writer)
	}
	return nil
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

func PrintCBOR(value interface{}, writer io.Writer) error {
	return WriteCBOR(value, writer)
}
