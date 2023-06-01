package streampackage

import (
	"archive/tar"
	contextpkg "context"
	"io"

	"github.com/klauspost/pgzip"
	"github.com/tliron/exturl"
)

//
// TarGZipStreamPackage
//

type TarGZipStreamPackage struct {
	reader     io.ReadCloser
	gzipReader *pgzip.Reader
	tarReader  *tar.Reader
}

func NewTarGZipStreamPackage(context contextpkg.Context, url exturl.URL) (*TarGZipStreamPackage, error) {
	var self TarGZipStreamPackage

	var err error
	if self.reader, err = url.Open(context); err == nil {
		if self.gzipReader, err = pgzip.NewReader(self.reader); err == nil {
			self.tarReader = tar.NewReader(self.gzipReader)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &self, nil
}

// StreamPackage interface
func (self *TarGZipStreamPackage) Next() (Stream, error) {
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
func (self *TarGZipStreamPackage) Close() error {
	self.gzipReader.Close() // TODO: err?
	return self.reader.Close()
}
