package streampackage

import (
	"archive/tar"
	contextpkg "context"
	"io"
	"io/fs"

	"github.com/tliron/exturl"
	"github.com/tliron/go-kutil/util"
)

//
// TarStreamPackage
//

type TarStreamPackage struct {
	reader    io.ReadCloser
	tarReader *tar.Reader
}

func NewTarStreamPackage(context contextpkg.Context, url exturl.URL) (*TarStreamPackage, error) {
	var self TarStreamPackage

	var err error
	if self.reader, err = url.Open(context); err == nil {
		self.tarReader = tar.NewReader(self.reader)
	} else {
		return nil, err
	}

	return &self, nil
}

// StreamPackage interface
func (self *TarStreamPackage) Next() (Stream, error) {
	for {
		header, err := self.tarReader.Next()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			} else {
				return nil, err
			}
		}
		if header.Typeflag == tar.TypeReg {
			return NewTarStream(header, self.tarReader), nil
		}
	}
}

// StreamPackage interface
func (self *TarStreamPackage) Close() error {
	return self.reader.Close()
}

//
// TarStream
//

type TarStream struct {
	header    *tar.Header
	tarReader *tar.Reader
}

func NewTarStream(header *tar.Header, tarReader *tar.Reader) *TarStream {
	return &TarStream{
		header:    header,
		tarReader: tarReader,
	}
}

// Stream interface
func (self *TarStream) Open(context contextpkg.Context) (io.ReadCloser, string, bool, error) {
	return io.NopCloser(self.tarReader), self.header.Name, util.IsFileExecutable(fs.FileMode(self.header.Mode)), nil
}
