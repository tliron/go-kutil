package terminal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zchee/color/v2"
)

var colorize = false

func EnableColor(force bool) {
	if force {
		color.NoColor = false
	}
	if color.NoColor {
		return
	}
	colorize = true
	Stdout = color.Output
	Stderr = color.Error
}

func ProcessColorizeFlag(colorize string) error {
	if colorize == "force" {
		EnableColor(true)
	} else if colorize_, err := strconv.ParseBool(colorize); err == nil {
		if colorize_ {
			EnableColor(false)
		}
	} else {
		return fmt.Errorf("\"--colorize\" must be boolean or \"force\": %s", colorize)
	}
	return nil
}

type Colorizer = func(name string) string

func ColorHeading(name string) string {
	if colorize {
		return color.GreenString(strings.ToUpper(name))
	} else {
		return name
	}
}

func ColorPath(name string) string {
	if colorize {
		return color.CyanString(name)
	} else {
		return name
	}
}

func ColorName(name string) string {
	if colorize {
		return color.BlueString(name)
	} else {
		return name
	}
}

func ColorTypeName(name string) string {
	if colorize {
		return color.MagentaString(name)
	} else {
		return name
	}
}

func ColorValue(name string) string {
	if colorize {
		return color.YellowString(name)
	} else {
		return name
	}
}

func ColorError(name string) string {
	if colorize {
		return color.RedString(name)
	} else {
		return name
	}
}
