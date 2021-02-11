package ard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v3"
)

func Read(reader io.Reader, format string, locate bool) (Map, Locator, error) {
	switch format {
	case "yaml", "":
		return ReadYAML(reader, locate)
	case "json":
		return ReadJSON(reader, locate)
	case "cjson":
		return ReadCompatibleJSON(reader, locate)
	case "xml":
		return ReadCompatibleXML(reader, locate)
	default:
		return nil, nil, fmt.Errorf("unsupported format: %q", format)
	}
}

func ReadYAML(reader io.Reader, locate bool) (Map, Locator, error) {
	var data Map
	var locator Locator
	var node yaml.Node

	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&node); err == nil {
		if decoded, err := yamlkeys.DecodeNode(&node); err == nil {
			var ok bool
			if data, ok = decoded.(Map); ok {
				if locate {
					locator = NewYAMLLocator(&node)
				}
			} else {
				return nil, nil, fmt.Errorf("YAML content is a \"%T\" instead of a map", decoded)
			}
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, yamlkeys.WrapWithDecodeError(err)
	}

	// We do not need to call EnsureMaps because yamlkeys takes care of it
	return data, locator, nil
}

func ReadJSON(reader io.Reader, locate bool) (Map, Locator, error) {
	var data StringMap
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		return EnsureMaps(data), nil, nil
	} else {
		return nil, nil, err
	}
}

func ReadCompatibleJSON(reader io.Reader, locate bool) (Map, Locator, error) {
	var data StringMap
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		map_ := FromCompatibleJSON(data)
		if map__, ok := map_.(Map); ok {
			return map__, nil, nil
		} else {
			return nil, nil, fmt.Errorf("JSON content is a \"%T\" instead of a map", map_)
		}
	} else {
		return nil, nil, err
	}
}

func ReadCompatibleXML(reader io.Reader, locate bool) (Map, Locator, error) {
	document := etree.NewDocument()
	if _, err := document.ReadFrom(reader); err == nil {
		elements := document.ChildElements()
		if len(elements) == 1 {
			if map_, err := FromCompatibleXML(elements[0]); err == nil {
				if map__, ok := map_.(Map); ok {
					return map__, nil, err
				} else {
					return nil, nil, errors.New("unsupported XML")
				}
			} else {
				return nil, nil, err
			}
		} else {
			return nil, nil, errors.New("unsupported XML")
		}
	} else {
		return nil, nil, err
	}
}
