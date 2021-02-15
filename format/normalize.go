package format

import (
	"strings"

	"github.com/tliron/yamlkeys"
)

func Normalize(value interface{}) (interface{}, error) {
	// TODO: not very efficient
	if code, err := EncodeYAML(value, " ", false); err == nil {
		return yamlkeys.Decode(strings.NewReader(code))
	} else {
		return nil, err
	}
}
