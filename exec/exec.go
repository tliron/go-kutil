package exec

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/creack/pty"
	"github.com/tliron/kutil/util"
)

// Caller has to close stdin, otherwise there will be a goroutine leak!
func ExecInteractive(pseudoTerminal bool, done chan error, name string, args ...string) (chan struct{}, chan []byte, chan []byte, chan []byte, error) {
	kill := make(chan struct{})
	stdin := make(chan []byte)
	stdout := make(chan []byte)
	stderr := make(chan []byte)

	command := exec.Command(name, args...)

	start := func(stdinWriter io.WriteCloser, stdoutReader io.Reader) {
		// Write stdout
		if stdoutReader != nil {
			go func() {
				if _, err := io.Copy(util.NewChannelWriter(stdout), stdoutReader); err != nil {
					// When done will return an input/output error
					log.Debugf("stdout copy error: %s", err.Error())
				}
			}()
		}

		// Read stdin, kill
		go func() {
			for {
				select {
				case b := <-stdin:
					if b == nil {
						log.Info("stdin closed")
						return
					}
					if _, err := stdinWriter.Write(b); err != nil {
						log.Errorf("stdin error: %s", err.Error())
						return
					}

				case <-kill:
					log.Info("killing process")
					if err := command.Process.Kill(); err != nil {
						log.Errorf("kill error: %s", err.Error())
					}
					return
				}
			}
		}()

		// Wait for command to end
		go func() {
			var doneError error
			if err := command.Wait(); err == nil {
				log.Info("command exited")
			} else {
				if exitError, ok := err.(*exec.ExitError); ok {
					doneError = fmt.Errorf("%d, %s", exitError.ExitCode(), exitError.Stderr)
				} else {
					log.Errorf("command wait error: %s", err.Error())
				}
			}

			if err := stdinWriter.Close(); err != nil {
				log.Errorf("stdin error: %s", err.Error())
			}
			close(stdout)
			close(stderr)

			done <- doneError
		}()
	}

	if pseudoTerminal {
		// Note: Our stderr is not a TTY file, which may cause some shell programs to disable their interactive mode.
		// In such cases it may be possible to force interactive mode, for example: `bash -i`
		command.Stderr = util.NewChannelWriter(stderr)
		if file, err := pty.Start(command); err == nil {
			start(file, file)
			return kill, stdin, stdout, stderr, nil
		} else {
			return nil, nil, nil, nil, err
		}
	} else {
		if stdinWriter, err := command.StdinPipe(); err == nil {
			command.Stdout = util.NewChannelWriter(stdout)
			command.Stderr = util.NewChannelWriter(stderr)
			if err := command.Start(); err == nil {
				start(stdinWriter, nil)
				return kill, stdin, stdout, stderr, nil
			} else {
				return nil, nil, nil, nil, err
			}
		} else {
			return nil, nil, nil, nil, err
		}
	}
}
