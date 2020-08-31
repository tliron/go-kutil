package format

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	templatepkg "text/template"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v3"
)

var YAMLNodeKinds = map[yaml.Kind]string{
	yaml.DocumentNode: "Document",
	yaml.SequenceNode: "Sequence",
	yaml.MappingNode:  "Mapping",
	yaml.ScalarNode:   "Scalar",
	yaml.AliasNode:    "Alias",
}

func WriteYAMLNodes(writer io.Writer, node *yaml.Node) {
	WriteYAMLNode(writer, node, 0)
}

func WriteYAMLNode(writer io.Writer, node *yaml.Node, indent int) {
	s := ""

	s += strings.Repeat(" ", indent)

	s += YAMLNodeKinds[node.Kind]

	switch node.Kind {
	// Document and alias tag is always "", nothing to print
	// Sequence tag is always "!!seq", no need to print
	// Mapping tag is always "!!map", no need to print

	case yaml.ScalarNode:
		s += " "
		s += node.Tag
	}

	if node.Value != "" {
		s += " "
		s += node.Value
	}

	fmt.Fprintln(writer, s)

	indent += 1
	for _, child := range node.Content {
		WriteYAMLNode(writer, child, indent)
	}
}

func DecodeYAMLStringMaps(code string) ([]ard.StringMap, error) {
	if values, err := yamlkeys.DecodeStringAll(code); err == nil {
		var r []ard.StringMap
		for _, value := range values {
			r = append(r, ard.EnsureStringMaps(value))
		}
		return r, nil
	} else {
		return nil, err
	}
}

func DecodeYAMLTemplate(code string, data interface{}) (ard.StringMap, error) {
	if template, err := templatepkg.New("").Parse(code); err == nil {
		var buffer bytes.Buffer
		if err := template.Execute(&buffer, data); err == nil {
			if value, err := yamlkeys.DecodeString(buffer.String()); err == nil {
				return ard.EnsureStringMaps(value), nil
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
