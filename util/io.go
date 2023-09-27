package util

import (
	contextpkg "context"
	"io"
	"sync"
	"testing"
)

const BUFFER_SIZE = 65536

func WriteNewline(writer io.Writer) error {
	_, err := io.WriteString(writer, "\n")
	return err
}

func ReaderSize(reader io.Reader) (int64, error) {
	var size int64 = 0

	buffer := make([]byte, BUFFER_SIZE)
	for {
		count, err := reader.Read(buffer)
		size += int64(count)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}
	}

	return size, nil
}

func ContextualRead(context contextpkg.Context, reader io.Reader, p []byte) (int, error) {
	if err := context.Err(); err == nil {
		return reader.Read(p)
	} else {
		return 0, err
	}
}

//
// BufferedWriter
//

// https://www.reddit.com/r/golang/comments/6irpt1/is_there_a_golang_logging_library_around_that/dj8yicz?utm_source=share&utm_medium=web2x&context=3
// https://gobyexample.com/closing-channels

type BufferedWriter struct {
	writer io.Writer
	jobs   chan []byte
	close  chan struct{}
	closed chan struct{}
}

func NewBufferedWriter(writer io.Writer, size int) BufferedWriter {
	self := BufferedWriter{
		writer: writer,
		jobs:   make(chan []byte, size),
		close:  make(chan struct{}, 1),
		closed: make(chan struct{}, 1),
	}

	go self.run()

	return self
}

// ([io.Writer] interface)
func (self BufferedWriter) Write(p []byte) (int, error) {
	defer func() {
		if recover() != nil {
			// The channel was closed
			//fmt.Println("closed!!!!")
			self.writer.Write(p)
		}
	}()

	self.jobs <- p
	return len(p), nil
}

// ([io.Closer] interface)
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
				self.closed <- struct{}{}
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

// ([io.Writer] interface)
func (self *SyncedWriter) Write(p []byte) (int, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.Writer.Write(p)
}

// ([io.Closer] interface)
func (self *SyncedWriter) Close() error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if closer, ok := self.Writer.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}

//
// ChannelWriter
//

type ChannelWriter struct {
	ch chan []byte
}

func NewChannelWriter(ch chan []byte) *ChannelWriter {
	return &ChannelWriter{ch}
}

// ([io.Writer] interface)
func (self *ChannelWriter) Write(p []byte) (int, error) {
	if p != nil {
		// We are copying the slice because contents might change while sending to the channel
		p_ := append(p[:0:0], p...)
		self.ch <- p_
	}
	return len(p), nil
}

//
// ChannelReader
//

type ChannelReader struct {
	reader *io.PipeReader
}

func NewChannelReader(ch chan []byte) *ChannelReader {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		for p := range ch {
			if _, err := writer.Write(p); err != nil {
				return
			}
		}
	}()

	return &ChannelReader{reader: reader}
}

// ([io.Reader] interface)
func (self *ChannelReader) Read(p []byte) (n int, err error) {
	return self.reader.Read(p)
}

//
// ContextualReader
//

// https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer.html
// https://medium.com/@zombiezen/canceling-i-o-in-go-capn-proto-5ae8c09c5b29
// https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/
// https://github.com/muesli/cancelreader

type ContextualReader struct {
	reader  io.Reader
	context contextpkg.Context
}

func NewContextualReader(context contextpkg.Context, reader io.Reader) io.Reader {
	return &ContextualReader{reader: reader, context: context}
}

// ([io.Reader] interface)
func (self *ContextualReader) Read(p []byte) (int, error) {
	return ContextualRead(self.context, self.reader, p)
}

//
// ContextualReadCloser
//

type ContextualReadCloser struct {
	reader  io.ReadCloser
	context contextpkg.Context
}

func NewContextualReadCloser(context contextpkg.Context, reader io.ReadCloser) io.ReadCloser {
	return &ContextualReadCloser{reader: reader, context: context}
}

// ([io.Reader] interface)
func (self *ContextualReadCloser) Read(p []byte) (int, error) {
	return ContextualRead(self.context, self.reader, p)
}

// ([io.Closer] interface)
func (self ContextualReadCloser) Close() error {
	return self.reader.Close()
}

//
// TestLogWriter
//

type TestLogWriter struct {
	t *testing.T
}

func NewTestLogWriter(t *testing.T) *TestLogWriter {
	return &TestLogWriter{t}
}

// ([io.Writer] interface)
func (self *TestLogWriter) Write(p []byte) (n int, err error) {
	self.t.Helper()
	self.t.Log(p)
	return len(p), nil
}
