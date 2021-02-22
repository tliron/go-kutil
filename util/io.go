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
// BufferWriter
//

// https://www.reddit.com/r/golang/comments/6irpt1/is_there_a_golang_logging_library_around_that/dj8yicz?utm_source=share&utm_medium=web2x&context=3

// Note: already implements io.Closer interface
type BufferWriter chan []byte

// io.Writer interface
func (self BufferWriter) Write(p []byte) (int, error) {
	self <- p
	return len(p), nil
}

func NewBufferWriter(writer io.Writer, n int) BufferWriter {
	bufferWriter := make(BufferWriter, n)

	go func() {
		for p := range bufferWriter {
			writer.Write(p)
		}
	}()

	return bufferWriter
}

//
// SyncWriter
//

type SyncWriter struct {
	Writer io.Writer

	lock sync.Mutex
}

func NewSyncWriter(writer io.Writer) *SyncWriter {
	return &SyncWriter{
		Writer: writer,
	}
}

// io.Writer interface
func (self *SyncWriter) Write(p []byte) (int, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.Writer.Write(p)
}

// io.Closer interface
func (self *SyncWriter) Close() error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if closer, ok := self.Writer.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}
