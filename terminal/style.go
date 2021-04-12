package terminal

import (
	"strings"
)

//
// Stylist
//

type Stylist struct {
	Colorize bool
}

func NewStylist(colorize bool) *Stylist {
	return &Stylist{colorize}
}

func (self *Stylist) Heading(name string) string {
	if self.Colorize {
		return ColorGreen(strings.ToUpper(name))
	} else {
		return strings.ToUpper(name)
	}
}

func (self *Stylist) Path(name string) string {
	if self.Colorize {
		return ColorCyan(name)
	} else {
		return name
	}
}

func (self *Stylist) Name(name string) string {
	if self.Colorize {
		return ColorBlue(name)
	} else {
		return name
	}
}

func (self *Stylist) TypeName(name string) string {
	if self.Colorize {
		return ColorMagenta(name)
	} else {
		return name
	}
}

func (self *Stylist) Value(name string) string {
	if self.Colorize {
		return ColorYellow(name)
	} else {
		return name
	}
}

func (self *Stylist) Error(name string) string {
	if self.Colorize {
		return ColorRed(name)
	} else {
		return name
	}
}
