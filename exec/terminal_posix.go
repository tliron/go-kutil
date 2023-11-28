//go:build !windows && !wasm

package exec

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func NewTerminal() (*Terminal, error) {
	// See: https://stackoverflow.com/a/54423725
	/*exec.Command("/usr/bin/stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("/usr/bin/stty", "-F", "/dev/tty", "-echo").Run()
	util.OnExit(func() {
		exec.Command("/usr/bin/stty", "-F", "/dev/tty", "echo").Run()
	})*/

	self := Terminal{
		Resize: make(chan Size),
	}

	var err error
	stdin := int(os.Stdin.Fd())
	if self.termState, err = term.MakeRaw(stdin); err == nil {
		if width, height, err := term.GetSize(stdin); err == nil {
			self.InitialSize = &Size{Width: uint(width), Height: uint(height)}
		}

		self.sigwinch = make(chan os.Signal)
		signal.Notify(self.sigwinch, syscall.SIGWINCH)
		go func() {
			for range self.sigwinch {
				if width, height, err := term.GetSize(stdin); err == nil {
					self.Resize <- Size{uint(width), uint(height)}
				}
			}
			log.Debug("closed sigwinch")
		}()

		return &self, nil
	} else {
		return nil, err
	}
}

func (self *Terminal) Close() error {
	os.Stdout.WriteString("\r")
	signal.Stop(self.sigwinch)
	close(self.sigwinch)
	close(self.Resize)
	stdin := int(os.Stdin.Fd())
	return term.Restore(stdin, self.termState)
}
