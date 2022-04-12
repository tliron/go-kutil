package util

import (
	"io"
	"os"
	"syscall"
)

//
// NamedPipe
//

type NamedPipe struct {
	path   string
	writer *os.File
	reader *os.File
}

func NewTempNamedPipe(writer io.Writer, mode uint32) (*NamedPipe, error) {
	var self NamedPipe

	pipe, err := os.CreateTemp("", "kutil-named-pipe-")
	if err != nil {
		return nil, err
	}
	self.path = pipe.Name()

	err = os.Remove(self.path)
	if err != nil {
		return nil, err
	}

	err = syscall.Mkfifo(self.path, mode)
	if err != nil {
		return nil, err
	}

	self.writer, err = os.OpenFile(self.path, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		os.Remove(self.path)
		return nil, err
	}

	self.reader, err = os.OpenFile(self.path, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		self.writer.Close()
		os.Remove(self.path)
		return nil, err
	}

	go io.Copy(writer, self.reader)

	return &self, nil
}

// io.Writer interface
func (self *NamedPipe) Write(p []byte) (int, error) {
	return self.writer.Write(p)
}

// io.Closer interface
func (self *NamedPipe) Close() error {
	self.reader.Close()
	self.writer.Close()
	return os.Remove(self.path)
}
