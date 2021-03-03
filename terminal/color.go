package terminal

import (
	"fmt"
	"strconv"

	"github.com/zchee/color/v2"
)

var Colorize = false

func EnableColor(force bool) {
	if force {
		color.NoColor = false
	}
	if color.NoColor {
		return
	}
	Colorize = true
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

func ColorRed(name string) string {
	if Colorize {
		return color.RedString(name)
	} else {
		return name
	}
}

func ColorGreen(name string) string {
	if Colorize {
		return color.GreenString(name)
	} else {
		return name
	}
}

func ColorYellow(name string) string {
	if Colorize {
		return color.YellowString(name)
	} else {
		return name
	}
}

func ColorBlue(name string) string {
	if Colorize {
		return color.BlueString(name)
	} else {
		return name
	}
}

func ColorMagenta(name string) string {
	if Colorize {
		return color.MagentaString(name)
	} else {
		return name
	}
}

func ColorCyan(name string) string {
	if Colorize {
		return color.CyanString(name)
	} else {
		return name
	}
}
