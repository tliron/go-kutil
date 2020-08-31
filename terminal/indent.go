package terminal

import (
	"fmt"
	"strings"
)

const IndentSpaces = 2

var Indent = strings.Repeat(" ", IndentSpaces)

func IndentString(indent int) string {
	return strings.Repeat(Indent, indent)
}

func PrintIndent(indent int) {
	fmt.Fprint(Stdout, IndentString(indent))
}
