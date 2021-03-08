package util

import (
	"io"
	"sync"
)

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
	jobs chan []byte
	done chan bool
}

func NewBufferedWriter(writer io.Writer, n int) BufferedWriter {
	self := BufferedWriter{
		jobs: make(chan []byte, n),
		done: make(chan bool, 1),
	}

	go self.run(writer)

	return self
}

// io.Writer interface
func (self BufferedWriter) Write(p []byte) (int, error) {
	self.jobs <- p
	return len(p), nil
}

// io.Closer interface
func (self BufferedWriter) Close() error {
	close(self.jobs)
	<-self.done
	return nil
}

func (self BufferedWriter) run(writer io.Writer) {
	for {
		if job, more := <-self.jobs; more {
			writer.Write(job)
		} else {
			self.done <- true
			break
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
