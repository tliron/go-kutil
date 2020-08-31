package format

import (
	"fmt"
	"strings"

	"github.com/tliron/kutil/terminal"
)

func Encode(data interface{}, format string, strict bool) (string, error) {
	switch format {
	case "yaml", "":
		return EncodeYAML(data, terminal.Indent, strict)
	case "json":
		return EncodeJSON(data, terminal.Indent)
	case "xml":
		return EncodeXML(data, terminal.Indent)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func EncodeYAML(data interface{}, indent string, strict bool) (string, error) {
	var writer strings.Builder
	if err := WriteYAML(data, &writer, indent, strict); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func EncodeJSON(data interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteJSON(data, &writer, indent); err == nil {
		s := writer.String()
		if indent == "" {
			// json.Encoder adds a "\n", unlike json.Marshal
			s = strings.Trim(s, "\n")
		}
		return s, nil
	} else {
		return "", err
	}
}

func EncodeXML(data interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteXML(data, &writer, indent); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}
