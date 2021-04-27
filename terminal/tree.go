package terminal

//
// TreePrefix
//

type TreePrefix []bool

func (self TreePrefix) Print(indent int, last bool) {
	PrintIndent(indent)

	for _, element := range self {
		if element {
			Print(Stdout, "  ")
		} else {
			Print(Stdout, "│ ")
		}
	}

	if last {
		Print(Stdout, "└─")
	} else {
		Print(Stdout, "├─")
	}
}
