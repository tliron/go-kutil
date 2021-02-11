package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/yamlkeys"
)

func Read(reader io.Reader, format string) (ard.Value, error) {
	switch format {
	case "yaml", "":
		return ReadYAML(reader)
	case "json":
		return ReadJSON(reader)
	case "cjson":
		return ReadCompatibleJSON(reader)
	case "xml":
		return ReadCompatibleXML(reader)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func ReadYAML(reader io.Reader) (ard.Value, error) {
	return yamlkeys.Decode(reader)
}

func ReadAllYAML(reader io.Reader) (ard.List, error) {
	return yamlkeys.DecodeAll(reader)
}

func ReadJSON(reader io.Reader) (ard.Value, error) {
	var data ard.Value
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		return ard.EnsureMaps(data), nil
	} else {
		return nil, err
	}
}

func ReadCompatibleJSON(reader io.Reader) (ard.Value, error) {
	var data ard.Value
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		return ard.FromCompatibleJSON(data), nil
	} else {
		return nil, err
	}
}

func ReadCompatibleXML(reader io.Reader) (ard.Value, error) {
	document := etree.NewDocument()
	if _, err := document.ReadFrom(reader); err == nil {
		elements := document.ChildElements()
		if len(elements) == 1 {
			return ard.FromCompatibleXML(elements[0])
		} else {
			return nil, fmt.Errorf("unsupported XML structure")
		}
	} else {
		return nil, err
	}
}
