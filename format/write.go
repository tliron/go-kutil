package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/ard"
	"gopkg.in/yaml.v3"
)

func Write(value interface{}, format string, indent string, strict bool, writer io.Writer) error {
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

	case "go":
		return WriteGo(value, writer, indent)

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteYAML(value interface{}, writer io.Writer, indent string, strict bool) error {
	if strict {
		value = ard.ToYAMLDocumentNode(value, true)
	}

	encoder := yaml.NewEncoder(writer)

	encoder.SetIndent(len(indent)) // This might not work as expected for tabs!
	// BUG: currently does not allow an indent value of 1, see: https://github.com/go-yaml/yaml/issues/501

	if slice, ok := value.([]interface{}); !ok {
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

func WriteJSON(value interface{}, writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	return encoder.Encode(value)
}

func WriteCompatibleJSON(value interface{}, writer io.Writer, indent string) error {
	return WriteJSON(ard.EnsureCompatibleJSON(value), writer, indent)
}

func WriteCompatibleXML(value interface{}, writer io.Writer, indent string) error {
	// Because we don't provide explicit marshalling for XML in the codebase (as we do for
	// JSON and YAML) we must canonicalize the data before encoding it
	value, err := ard.Canonicalize(value)
	if err != nil {
		return err
	}

	value = ard.ToCompatibleXML(value)

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

func WriteCBOR(value interface{}, writer io.Writer) error {
	encoder := cbor.NewEncoder(writer)
	return encoder.Encode(value)
}

func WriteGo(value interface{}, writer io.Writer, indent string) error {
	NewUtterConfig(indent).Fdump(writer, value)
	return nil
}
