package streampackage

import (
	contextpkg "context"
	"io"

	"github.com/klauspost/compress/zip"
	"github.com/tliron/go-kutil/util"
)

//
// ZipStreamPackage
//

type ZipStreamPackage struct {
	zipReader *zip.ReadCloser

	streams []*ZipStream
	index   int
}

func NewZipStreamPackage(path string) (*ZipStreamPackage, error) {
	var self ZipStreamPackage
	var err error
	if self.zipReader, err = zip.OpenReader(path); err == nil {
		for _, file := range self.zipReader.File {
			if !file.FileInfo().IsDir() {
				self.streams = append(self.streams, NewZipStream(file))
			}
		}
		return &self, nil
	} else {
		return nil, err
	}
}

// StreamPackage interface
func (self *ZipStreamPackage) Next() (Stream, error) {
	if self.index < len(self.streams) {
		source := self.streams[self.index]
		self.index++
		return source, nil
	} else {
		return nil, nil
	}
}

// StreamPackage interface
func (self *ZipStreamPackage) Close() error {
	return self.zipReader.Close()
}

//
// ZipStream
//

type ZipStream struct {
	file *zip.File

	reader io.ReadCloser
}

func NewZipStream(file *zip.File) *ZipStream {
	return &ZipStream{
		file: file,
	}
}

// Stream interface
func (self *ZipStream) Open(context contextpkg.Context) (io.ReadCloser, string, bool, error) {
	var err error
	if self.reader, err = self.file.Open(); err == nil {
		return self.reader, self.file.Name, util.IsFileExecutable(self.file.Mode()), nil
	} else {
		return nil, "", false, err
	}
}
