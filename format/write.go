package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/kutil/ard"
	"gopkg.in/yaml.v3"
)

func Write(data interface{}, format string, indent string, strict bool, writer io.Writer) error {
	// Special handling for bare strings
	if s, ok := data.(string); ok {
		_, err := io.WriteString(writer, s)
		return err
	}

	// Special handling for XML document (etree)
	if xmlDocument, ok := data.(*etree.Document); ok {
		return WriteXMLDocument(xmlDocument, writer, indent)
	}

	switch format {
	case "yaml", "":
		return WriteYAML(data, writer, indent, strict)
	case "json":
		return WriteJSON(data, writer, indent)
	case "cjson":
		return WriteCompatibleJSON(data, writer, indent)
	case "xml":
		return WriteXML(data, writer, indent)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteYAML(data interface{}, writer io.Writer, indent string, strict bool) error {
	if strict {
		data = ard.ToYAMLDocumentNode(data, true)
	}

	encoder := yaml.NewEncoder(writer)

	encoder.SetIndent(len(indent)) // This might not work as expected for tabs!
	// BUG: currently does not allow an indent value of 1, see: https://github.com/go-yaml/yaml/issues/501

	if slice, ok := data.([]interface{}); !ok {
		return encoder.Encode(data)
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

func WriteJSON(data interface{}, writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	return encoder.Encode(data)
}

func WriteCompatibleJSON(data interface{}, writer io.Writer, indent string) error {
	return WriteJSON(ToCompatibleJSON(data), writer, indent)
}

func WriteXML(data interface{}, writer io.Writer, indent string) error {
	// Because we don't provide explicit marshalling for XML in the codebase (as we do for
	// JSON and YAML) we must normalize the data before encoding it
	data, err := Normalize(data)
	if err != nil {
		return err
	}

	data = ToXMLWritable(data)

	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", indent)
	if err := encoder.Encode(data); err != nil {
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
