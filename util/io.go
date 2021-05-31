package util

import (
	"fmt"
	"io"
	"sync"
)

func WriteNewline(writer io.Writer) error {
	_, err := io.WriteString(writer, "\n")
	return err
}

func ReaderSize(reader io.Reader) (int64, error) {
	var size int64 = 0

	buffer := make([]byte, 1024)
	for {
		if count, err := reader.Read(buffer); err == nil {
			size += int64(count)
		} else if err == io.EOF {
			break
		} else {
			return 0, err
		}
	}

	return size, nil
}

//
// BufferedWriter
//

// https://www.reddit.com/r/golang/comments/6irpt1/is_there_a_golang_logging_library_around_that/dj8yicz?utm_source=share&utm_medium=web2x&context=3
// https://gobyexample.com/closing-channels

type BufferedWriter struct {
	writer io.Writer
	jobs   chan []byte
	close  chan bool
	closed chan bool
}

func NewBufferedWriter(writer io.Writer, size int) BufferedWriter {
	self := BufferedWriter{
		writer: writer,
		jobs:   make(chan []byte, size),
		close:  make(chan bool, 1),
		closed: make(chan bool, 1),
	}

	go self.run()

	return self
}

func (self BufferedWriter) CloseOnExit() ExitFunctionHandle {
	return OnExit(func() {
		self.Close()
	})
}

// io.Writer interface
func (self BufferedWriter) Write(p []byte) (int, error) {
	defer func() {
		if recover() != nil {
			// The channel was closed
			fmt.Println("closed!!!!")
			self.writer.Write(p)
		}
	}()

	self.jobs <- p
	return len(p), nil
}

// io.Closer interface
func (self BufferedWriter) Close() error {
	defer func() {
		recover()
	}()

	close(self.jobs)
	<-self.closed
	return nil
}

func (self BufferedWriter) run() {
	for {
		select {
		case job, ok := <-self.jobs:
			if ok {
				self.writer.Write(job)
			} else {
				self.closed <- true
				return
			}
		}
	}
}

//
// SyncedWriter
//

type SyncedWriter struct {
	Writer io.Writer

	lock sync.Mutex
}

func NewSyncedWriter(writer io.Writer) *SyncedWriter {
	return &SyncedWriter{
		Writer: writer,
	}
}

// io.Writer interface
func (self *SyncedWriter) Write(p []byte) (int, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.Writer.Write(p)
}

// io.Closer interface
func (self *SyncedWriter) Close() error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if closer, ok := self.Writer.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}
