package kubernetes

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restpkg "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func WriteToContainer(rest restpkg.Interface, config *restpkg.Config, namespace string, podName string, containerName string, reader io.Reader, targetPath string, permissions *int64) error {
	dir := filepath.Dir(targetPath)
	if err := Exec(rest, config, namespace, podName, containerName, nil, nil, nil, false, "mkdir", "--parents", dir); err == nil {
		if err := Exec(rest, config, namespace, podName, containerName, reader, nil, nil, false, "cp", "/dev/stdin", targetPath); err == nil {
			if permissions != nil {
				octal := strconv.FormatInt(*permissions, 8)
				return Exec(rest, config, namespace, podName, containerName, nil, nil, nil, false, "chmod", octal, targetPath)
			} else {
				return nil
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func ReadFromContainer(rest restpkg.Interface, config *restpkg.Config, namespace string, podName string, containerName string, writer io.Writer, sourcePath string) error {
	return Exec(rest, config, namespace, podName, containerName, nil, writer, nil, false, "cat", sourcePath)
}

func Exec(rest restpkg.Interface, config *restpkg.Config, namespace string, podName string, containerName string, stdin io.Reader, stdout io.Writer, stderr io.Writer, tty bool, command ...string) error {
	var stderrCapture strings.Builder
	if stderr == nil {
		// If not redirecting stderr then make sure to capture it
		stderr = &stderrCapture
	}

	execOptions := core.PodExecOptions{
		Container: containerName,
		Command:   command,
		TTY:       tty,
		Stderr:    true,
	}

	streamOptions := remotecommand.StreamOptions{
		Tty:    tty,
		Stderr: stderr,
	}

	if stdin != nil {
		execOptions.Stdin = true
		streamOptions.Stdin = stdin
	}

	if stdout != nil {
		execOptions.Stdout = true
		streamOptions.Stdout = stdout
	}

	request := rest.Post().Namespace(namespace).Resource("pods").Name(podName).SubResource("exec").VersionedParams(&execOptions, scheme.ParameterCodec)

	if executor, err := remotecommand.NewSPDYExecutor(config, "POST", request.URL()); err == nil {
		if err = executor.Stream(streamOptions); err == nil {
			return nil
		} else {
			return NewExecError(err, strings.TrimRight(stderrCapture.String(), "\n"))
		}
	} else {
		return err
	}
}

type ExecError struct {
	Err    error
	Stderr string
}

func NewExecError(err error, stderr string) *ExecError {
	return &ExecError{err, stderr}
}

// (error interface)
func (self *ExecError) Error() string {
	return fmt.Sprintf("%s\n%s", self.Err.Error(), self.Stderr)
}
