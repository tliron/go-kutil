package ard

import (
	"bytes"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v3"
)

func AgnosticCopy(value Value) (Value, error) {
	if IsPrimitiveType(value) {
		return value, nil
	} else {
		var err error
		switch value_ := value.(type) {
		case Map:
			map_ := make(Map)
			for key, value_ := range value_ {
				if map_[key], err = AgnosticCopy(value_); err != nil {
					return nil, err
				}
			}
			return map_, nil

		case StringMap:
			map_ := make(StringMap)
			for key, value_ := range value_ {
				if map_[key], err = AgnosticCopy(value_); err != nil {
					return nil, err
				}
			}
			return map_, nil

		case List:
			list := make(List, len(value_))
			for index, entry := range value_ {
				if list[index], err = AgnosticCopy(entry); err != nil {
					return nil, err
				}
			}
			return list, nil

		default:
			// TODO: not very efficient
			return AgnosticCopyThroughCBOR(value)
		}
	}
}

func AgnosticCopyThroughCBOR(value Value) (Value, error) {
	if code, err := cbor.Marshal(value); err == nil {
		var value_ Value
		if err := cbor.Unmarshal(code, &value_); err == nil {
			return value_, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func AgnosticCopyThroughYAML(value Value) (Value, error) {
	if code, err := yaml.Marshal(value); err == nil {
		if value, err := yamlkeys.Decode(bytes.NewReader(code)); err == nil {
			return value, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NormalizeMapsAgnosticCopy(value Value) (Value, error) {
	if value_, err := AgnosticCopy(value); err == nil {
		value_, _ = NormalizeMaps(value_)
		return value_, nil
	} else {
		return nil, err
	}
}

func NormalizeStringMapsAgnosticCopy(value Value) (Value, error) {
	if value_, err := AgnosticCopy(value); err == nil {
		value_, _ = NormalizeStringMaps(value_)
		return value_, nil
	} else {
		return nil, err
	}
}

func SimpleCopy(value Value) Value {
	switch value_ := value.(type) {
	case Map:
		map_ := make(Map)
		for key, value_ := range value_ {
			map_[key] = SimpleCopy(value_)
		}
		return map_

	case StringMap:
		map_ := make(StringMap)
		for key, value_ := range value_ {
			map_[key] = SimpleCopy(value_)
		}
		return map_

	case List:
		list := make(List, len(value_))
		for index, entry := range value_ {
			list[index] = SimpleCopy(entry)
		}
		return list

	default:
		return value
	}
}

func NormalizeMapsSimpleCopy(value Value) Value {
	value_ := SimpleCopy(value)
	value_, _ = NormalizeMaps(value_)
	return value_
}

func NormalizeStringMapsSimpleCopy(value Value) Value {
	value_ := SimpleCopy(value)
	value_, _ = NormalizeStringMaps(value_)
	return value_
}
