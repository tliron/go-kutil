package ard

import (
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/util"
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

	case "cbor":
		return DecodeCBOR(code, locate)

	default:
		return nil, nil, fmt.Errorf("unsupported format: %q", format)
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

// The code should be in Base64
func DecodeCBOR(code string, locate bool) (Value, Locator, error) {
	var value Value
	if bytes, err := util.FromBase64(code); err == nil {
		if err := cbor.Unmarshal(bytes, &value); err == nil {
			return value, nil, nil
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}
}
