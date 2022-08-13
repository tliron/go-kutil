package transcribe

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
	"gopkg.in/yaml.v3"
)

func Write(value any, format string, indent string, strict bool, writer io.Writer) error {
	// Special handling for bare strings (format is ignored)
	if s, ok := value.(string); ok {
		_, err := io.WriteString(writer, s)
		return err
	}

	// Special handling for XML etree document (format is ignored)
	if xmlDocument, ok := value.(*etree.Document); ok {
		return WriteXMLDocument(xmlDocument, writer, indent)
	}

	switch format {
	case "yaml", "":
		return WriteYAML(value, writer, indent, strict)

	case "json":
		return WriteJSON(value, writer, indent)

	case "cjson":
		return WriteCompatibleJSON(value, writer, indent)

	case "xml":
		return WriteCompatibleXML(value, writer, indent)

	case "cbor":
		return WriteCBOR(value, writer)

	case "messagepack":
		return WriteMessagePack(value, writer)

	case "go":
		return WriteGo(value, writer, indent)

	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func WriteYAML(value any, writer io.Writer, indent string, strict bool) error {
	if strict {
		value = ard.ToYAMLDocumentNode(value, true)
	}

	encoder := yaml.NewEncoder(writer)

	encoder.SetIndent(len(indent)) // This might not work as expected for tabs!
	// BUG: currently does not allow an indent value of 1, see: https://github.com/go-yaml/yaml/issues/501

	if slice, ok := value.([]any); !ok {
		return encoder.Encode(value)
	} else {
		// YAML separates each entry with "---"
		// (In JSON the slice would be written as an array)
		for _, data_ := range slice {
			if err := encoder.Encode(data_); err != nil {
				return err
			}
		}
		return nil
	}
}

func WriteJSON(value any, writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	return encoder.Encode(value)
}

func WriteCompatibleJSON(value any, writer io.Writer, indent string) error {
	if value_, err := ard.EnsureCompatibleJSON(value); err == nil {
		return WriteJSON(value_, writer, indent)
	} else {
		return err
	}
}

func WriteCompatibleXML(value any, writer io.Writer, indent string) error {
	value, err := ard.EnsureCompatibleXML(value)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", indent)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	if indent == "" {
		// When there's no indent the XML encoder does not emit a final newline
		// (We want it for consistency with YAML and JSON)
		if _, err := io.WriteString(writer, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func WriteXMLDocument(xmlDocument *etree.Document, writer io.Writer, indent string) error {
	xmlDocument.Indent(len(indent))
	_, err := xmlDocument.WriteTo(writer)
	return err
}

func WriteCBOR(value any, writer io.Writer) error {
	encoder := cbor.NewEncoder(writer)
	return encoder.Encode(value)
}

func WriteMessagePack(value any, writer io.Writer) error {
	encoder := util.NewMessagePackEncoder(writer)
	return encoder.Encode(value)
}

func WriteGo(value any, writer io.Writer, indent string) error {
	NewUtterConfig(indent).Fdump(writer, value)
	return nil
}
