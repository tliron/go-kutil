package format

import (
	"fmt"
	"strings"

	"github.com/tliron/kutil/terminal"
)

func Encode(value interface{}, format string, strict bool) (string, error) {
	switch format {
	case "yaml", "":
		return EncodeYAML(value, terminal.Indent, strict)

	case "json":
		return EncodeJSON(value, terminal.Indent)

	case "cjson":
		return EncodeCompatibleJSON(value, terminal.Indent)

	case "xml":
		return EncodeCompatibleXML(value, terminal.Indent)

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func EncodeYAML(value interface{}, indent string, strict bool) (string, error) {
	var writer strings.Builder
	if err := WriteYAML(value, &writer, indent, strict); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func EncodeJSON(value interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteJSON(value, &writer, indent); err == nil {
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

func EncodeCompatibleJSON(value interface{}, indent string) (string, error) {
	return EncodeJSON(ToCompatibleJSON(value), indent)
}

func EncodeCompatibleXML(value interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteCompatibleXML(value, &writer, indent); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}
