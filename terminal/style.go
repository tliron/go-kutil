package terminal

import (
	"strings"

	"github.com/zchee/color/v2"
)

func StyleHeading(name string) string {
	if Colorize {
		return color.GreenString(strings.ToUpper(name))
	} else {
		return strings.ToUpper(name)
	}
}

func StylePath(name string) string {
	if Colorize {
		return color.CyanString(name)
	} else {
		return name
	}
}

func StyleName(name string) string {
	if Colorize {
		return color.BlueString(name)
	} else {
		return name
	}
}

func StyleTypeName(name string) string {
	if Colorize {
		return color.MagentaString(name)
	} else {
		return name
	}
}

func StyleValue(name string) string {
	if Colorize {
		return color.YellowString(name)
	} else {
		return name
	}
}

func StyleError(name string) string {
	if Colorize {
		return color.RedString(name)
	} else {
		return name
	}
}
