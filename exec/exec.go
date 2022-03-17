package exec

import (
	"fmt"
	"os/exec"

	"github.com/tliron/kutil/util"
)

const CHANNEL_SIZE = 10

// Caller has to close stdin, otherwise there will be a goroutine leak!
func ExecInteractive(done chan error, name string, args ...string) (chan struct{}, chan []byte, chan []byte, chan []byte, error) {
	kill := make(chan struct{})
	stdin := make(chan []byte, CHANNEL_SIZE)
	stdout := make(chan []byte, CHANNEL_SIZE)
	stderr := make(chan []byte, CHANNEL_SIZE)

	command := exec.Command(name, args...)
	if stdinWriter, err := command.StdinPipe(); err == nil {
		command.Stdout = util.NewChannelWriter(stdout)
		command.Stderr = util.NewChannelWriter(stderr)
		if err := command.Start(); err == nil {

			// Read channels
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

			// Wait for command
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

				close(stdout)
				close(stderr)

				done <- doneError
			}()

			return kill, stdin, stdout, stderr, nil
		} else {
			return nil, nil, nil, nil, err
		}
	} else {
		return nil, nil, nil, nil, err
	}
}
