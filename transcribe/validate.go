package transcribe

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/util"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

func Validate(code string, format string) error {
	switch format {
	case "yaml", "":
		return ValidateYAML(code)

	case "json", "cjson":
		return ValidateJSON(code)

	case "xml":
		return ValidateXML(code)

	case "cbor":
		return ValidateCBOR(code)

	case "messagepack":
		return ValidateMessagePack(code)

	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func ValidateYAML(code string) error {
	decoder := yaml.NewDecoder(strings.NewReader(code))
	// Note: decoder.Decode will only decode the first document it finds
	for {
		var node yaml.Node
		if err := decoder.Decode(&node); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}

func ValidateJSON(code string) error {
	decoder := json.NewDecoder(strings.NewReader(code))
	// Note: decoder.Decode will only decode the first element it finds
	for {
		if _, err := decoder.Token(); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}

func ValidateXML(code string) error {
	decoder := xml.NewDecoder(strings.NewReader(code))
	// Note: decoder.Decode will only decode the first element it finds
	for {
		if _, err := decoder.Token(); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}

func ValidateCBOR(code string) error {
	var value any
	if bytes, err := util.FromBase64(code); err == nil {
		if err := cbor.Unmarshal(bytes, &value); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func ValidateMessagePack(code string) error {
	var value any
	if bytes, err := util.FromBase64(code); err == nil {
		if err := msgpack.Unmarshal(bytes, &value); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
