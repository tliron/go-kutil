//go:build windows || wasm

package exec

import (
	"errors"
)

func NewTerminal() (*Terminal, error) {
	return nil, errors.New("Terminal is not supported on this platform")
}

func (self *Terminal) Close() error {
	return errors.New("Terminal is not supported on this platform")
}
