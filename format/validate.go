package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

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
	default:
		return fmt.Errorf("unsupported format: %s", format)
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
