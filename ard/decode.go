package ard

import (
	"fmt"
	"strings"
)

func Decode(code string, format string, locate bool) (Value, Locator, error) {
	switch format {
	case "yaml", "":
		return DecodeYAML(code, locate)

	case "json":
		return DecodeJSON(code, locate)

	case "cjson":
		return DecodeCompatibleJSON(code, locate)

	case "xml":
		return DecodeCompatibleXML(code, locate)

	default:
		return nil, nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func DecodeYAML(code string, locate bool) (Value, Locator, error) {
	return ReadYAML(strings.NewReader(code), locate)
}

func DecodeJSON(code string, locate bool) (Value, Locator, error) {
	return ReadJSON(strings.NewReader(code), locate)
}

func DecodeCompatibleJSON(code string, locate bool) (Value, Locator, error) {
	return ReadCompatibleJSON(strings.NewReader(code), locate)
}

func DecodeCompatibleXML(code string, locate bool) (Value, Locator, error) {
	return ReadCompatibleXML(strings.NewReader(code), locate)
}
