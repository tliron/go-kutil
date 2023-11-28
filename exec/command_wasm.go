//go:build wasm

package exec

import (
	contextpkg "context"
	"errors"
)

func (self *Command) Start(context contextpkg.Context) (*Process, error) {
	return nil, errors.New("Command.Start is not supported on this platform")
}
