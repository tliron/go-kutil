package ard

import (
	"encoding/json"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v2"
)

func RoundtripCBOR(value Value) (Value, error) {
	if bytes, err := cbor.Marshal(value); err == nil {
		var value Value
		if err := cbor.Unmarshal(bytes, &value); err == nil {
			return value, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func RoundtripYAML(value Value) (Value, error) {
	var writer strings.Builder
	encoder := yaml.NewEncoder(&writer)
	if err := encoder.Encode(value); err == nil {
		return yamlkeys.Decode(strings.NewReader(writer.String()))
	} else {
		return nil, err
	}
}

func RoundtripCompatibleJSON(value Value) (Value, error) {
	var writer strings.Builder
	encoder := json.NewEncoder(&writer)
	value = EnsureCompatibleJSON(value)
	if err := encoder.Encode(value); err == nil {
		var value_ Value
		decoder := json.NewDecoder(strings.NewReader(writer.String()))
		if err := decoder.Decode(&value_); err == nil {
			value_, _ = FromCompatibleJSON(value_)
			return value_, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
