package exec

import (
	"os"
	"path/filepath"
)

const DEFAULT_CHANNEL_SIZE = 10

//
// Command
//

type Command struct {
	Name           string
	Args           []string
	Dir            string
	Environment    map[string]string
	PseudoTerminal *Size
	ChannelSize    int

	done chan error
}

func NewCommand() *Command {
	return &Command{
		ChannelSize: DEFAULT_CHANNEL_SIZE,
		done:        make(chan error),
	}
}

func (self *Command) Wait() error {
	err := <-self.done
	return err
}

func (self *Command) Stop(err error) {
	self.done <- err
}

func (self *Command) AddPath(key string, path string) {
	if self.Environment == nil {
		self.Environment = map[string]string{key: path}
	} else {
		if value, ok := self.Environment[key]; ok {
			if value == "" {
				self.Environment[key] = path
			} else {
				for _, path_ := range filepath.SplitList(value) {
					if path_ == path {
						return
					}
				}
				self.Environment[key] = value + string(os.PathListSeparator) + path
			}
		} else {
			self.Environment[key] = path
		}
	}
}
