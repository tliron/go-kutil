package terminal

import (
	"github.com/muesli/termenv"
)

func enableColor() (Cleanup, error) {
	if mode, err := termenv.EnableWindowsANSIConsole(); err == nil {
		return func() error {
			return termenv.RestoreWindowsConsole(mode)
		}, nil
	} else {
		return nil, err
	}
}
