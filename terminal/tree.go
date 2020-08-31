package terminal

import (
	"fmt"
)

//
// TreePrefix
//

type TreePrefix []bool

func (self TreePrefix) Print(indent int, last bool) {
	PrintIndent(indent)

	for _, element := range self {
		if element {
			fmt.Fprint(Stdout, "  ")
		} else {
			fmt.Fprint(Stdout, "│ ")
		}
	}

	if last {
		fmt.Fprint(Stdout, "└─")
	} else {
		fmt.Fprint(Stdout, "├─")
	}
}
