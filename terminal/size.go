// +build !wasm

package terminal

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetSize() (int, int, error) {
	return terminal.GetSize(int(os.Stdout.Fd()))
}
