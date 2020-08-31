package format

import (
	"fmt"
	"strings"

	"github.com/tliron/kutil/ard"
)

func Decode(code string, format string, all bool) (ard.Value, error) {
	switch format {
	case "yaml", "":
		if all {
			return DecodeAllYAML(code)
		} else {
			return DecodeYAML(code)
		}
	case "json":
		return DecodeJSON(code)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func DecodeYAML(code string) (ard.Value, error) {
	return ReadYAML(strings.NewReader(code))
}

func DecodeAllYAML(code string) (ard.List, error) {
	return ReadAllYAML(strings.NewReader(code))
}

func DecodeJSON(code string) (ard.Value, error) {
	return ReadJSON(strings.NewReader(code))
}
