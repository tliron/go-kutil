package ard

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/util"
	"gopkg.in/yaml.v3"
)

// Default is CBOR
func Roundtrip(value Value, format string) (Value, error) {
	switch format {
	case "yaml":
		return RoundtripYAML(value)

	case "cjson":
		return RoundtripCompatibleJSON(value)

	case "xml":
		return RoundtripCompatibleXML(value)

	case "cbor", "":
		return RoundtripCBOR(value)

	case "messagepack":
		return RoundtripMessagePack(value)

	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

func RoundtripYAML(value Value) (Value, error) {
	var writer strings.Builder
	encoder := yaml.NewEncoder(&writer)
	if err := encoder.Encode(value); err == nil {
		value_, _, err := ReadYAML(strings.NewReader(writer.String()), false)
		return value_, err
	} else {
		return nil, err
	}
}

func RoundtripCompatibleJSON(value Value) (Value, error) {
	value = EnsureCompatibleJSON(value)
	var writer strings.Builder
	encoder := json.NewEncoder(&writer)
	if err := encoder.Encode(value); err == nil {
		value_, _, err := ReadCompatibleJSON(strings.NewReader(writer.String()), false)
		return value_, err
	} else {
		return nil, err
	}
}

func RoundtripCompatibleXML(value Value) (Value, error) {
	// Because we don't provide explicit marshalling for XML in the codebase (as we do for
	// JSON and YAML) we must canonicalize the data before encoding it
	if value, err := Canonicalize(value); err == nil {
		value = ToCompatibleXML(value)
		var writer strings.Builder
		if _, err := writer.WriteString(xml.Header); err == nil {
			encoder := xml.NewEncoder(&writer)
			encoder.Indent("", "")
			if err := encoder.Encode(value); err == nil {
				value_, _, err := ReadCompatibleXML(strings.NewReader(writer.String()), false)
				return value_, err
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func RoundtripCBOR(value Value) (Value, error) {
	if bytes, err := cbor.Marshal(value); err == nil {
		var value_ Value
		if err := cbor.Unmarshal(bytes, &value_); err == nil {
			return value_, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func RoundtripMessagePack(value Value) (Value, error) {
	var buffer bytes.Buffer
	encoder := util.NewMessagePackEncoder(&buffer)
	if err := encoder.Encode(value); err == nil {
		value_, _, err := ReadMessagePack(&buffer)
		return value_, err
	} else {
		return nil, err
	}
}
