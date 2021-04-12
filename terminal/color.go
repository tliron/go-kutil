package terminal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zchee/color/v2"
)

var Colorize = false

var Stylize = NewStylist(false)

func EnableColor(force bool) {
	if force {
		color.NoColor = false
	}

	Colorize = !color.NoColor

	if Colorize {
		Stdout = color.Output
		Stderr = color.Error
		Stylize = NewStylist(true)
	} else {
		Stdout = os.Stdout
		Stderr = os.Stderr
		Stylize = NewStylist(false)
	}
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

// Colorizer signature
func ColorRed(name string) string {
	return color.RedString(name)
}

// Colorizer signature
func ColorGreen(name string) string {
	return color.GreenString(name)
}

// Colorizer signature
func ColorYellow(name string) string {
	return color.YellowString(name)
}

// Colorizer signature
func ColorBlue(name string) string {
	return color.BlueString(name)
}

// Colorizer signature
func ColorMagenta(name string) string {
	return color.MagentaString(name)
}

// Colorizer signature
func ColorCyan(name string) string {
	return color.CyanString(name)
}
