package transcribe

import (
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

var YAMLNodeKinds = map[yaml.Kind]string{
	yaml.DocumentNode: "Document",
	yaml.SequenceNode: "Sequence",
	yaml.MappingNode:  "Mapping",
	yaml.ScalarNode:   "Scalar",
	yaml.AliasNode:    "Alias",
}

func DumpYAMLNodes(writer io.Writer, node *yaml.Node) {
	DumpYAMLNode(writer, node, 0)
}

func DumpYAMLNode(writer io.Writer, node *yaml.Node, indent int) {
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

	io.WriteString(writer, s)

	indent += 1
	for _, child := range node.Content {
		DumpYAMLNode(writer, child, indent)
	}
}
