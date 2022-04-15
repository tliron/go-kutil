package exec

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/tliron/kutil/util"
)

const DEFAULT_CHANNEL_SIZE = 10

//
// Command
//

type Command struct {
	Name           string
	Args           []string
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

func (self *Command) Start() (*Process, error) {
	process := newProcess(self.ChannelSize)

	command := exec.Command(self.Name, self.Args...)
	if self.Environment != nil {
		for k, v := range self.Environment {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	start := func(stdinWriter io.WriteCloser, stdoutReader io.Reader, tty *os.File) {
		// Write stdout
		if stdoutReader != nil {
			go func() {
				io.Copy(util.NewChannelWriter(process.Stdout), stdoutReader)
				// When done will return an input/output error
				// (an *fs.PathError wrapping a syscall.Errno)
				log.Debug("stdout closed")
			}()
		}

		// Read kill, size, stdin
		go func() {
			for {
				select {
				case b := <-process.stdin:
					if b == nil {
						log.Debug("stdin closed")
						return
					}
					if _, err := stdinWriter.Write(b); err != nil {
						log.Errorf("stdin: %s", err.Error())
						return
					}

				case s := <-process.resize:
					if tty != nil {
						winsize := pty.Winsize{Rows: uint16(s.Height), Cols: uint16(s.Width)}
						if err := pty.Setsize(tty, &winsize); err != nil {
							log.Errorf("resize: %s", err.Error())
						}
					}

				case <-process.kill:
					log.Info("killing process")
					if err := command.Process.Kill(); err != nil {
						log.Errorf("kill: %s", err.Error())
					}
					return
				}
			}
		}()

		// Wait for command to end
		go func() {
			var exitError error

			if err := command.Wait(); err == nil {
				log.Info("command exited")
			} else {
				if err_, ok := err.(*exec.ExitError); ok {
					exitError = fmt.Errorf("%d: %s", err_.ExitCode(), err_.Stderr)
				} else {
					log.Errorf("command wait: %s", err.Error())
				}
			}

			if err := stdinWriter.Close(); err != nil {
				log.Errorf("stdin: %s", err.Error())
			}

			close(process.Stdout)
			close(process.Stderr)

			self.Stop(exitError)
		}()
	}

	log.Debugf("%s", command.String())

	if self.PseudoTerminal != nil {
		// Note: Our stderr is not a TTY file, which may cause some shell programs to disable their interactive mode.
		// In such cases it may be possible to force interactive mode, for example: `bash -i`
		log.Debugf("creating pseudo-terminal with size %d, %d", self.PseudoTerminal.Width, self.PseudoTerminal.Height)
		command.Stderr = util.NewChannelWriter(process.Stderr)
		winsize := pty.Winsize{Rows: uint16(self.PseudoTerminal.Height), Cols: uint16(self.PseudoTerminal.Width)}
		if ptyFile, err := pty.StartWithSize(command, &winsize); err == nil {
			start(ptyFile, ptyFile, ptyFile)
			return &process, nil
		} else {
			return nil, err
		}
	} else {
		if stdinWriter, err := command.StdinPipe(); err == nil {
			command.Stdout = util.NewChannelWriter(process.Stdout)
			command.Stderr = util.NewChannelWriter(process.Stderr)
			if err := command.Start(); err == nil {
				start(stdinWriter, nil, nil)
				return &process, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
}
