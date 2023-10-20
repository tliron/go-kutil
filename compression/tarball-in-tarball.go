package compression

import (
	"archive/tar"
	"io"
	"strings"
	"sync"

	"github.com/klauspost/pgzip"
)

func OpenFirstTarballInTarball(reader io.Reader) (io.Reader, error) {
	tarReader := tar.NewReader(reader)

	for {
		// TODO: support more compression types
		if header, err := tarReader.Next(); err == nil {
			if (header.Typeflag == tar.TypeReg) && strings.HasSuffix(header.Name, ".tar.gz") {
				return pgzip.NewReader(tarReader)
			}
		} else {
			return nil, err // can be io.EOF
		}
	}
}

//
// FirstTarballInTarballDecoder
//
// Decodes the first tar entry with a ".tar.gz" extension
//

type FirstTarballInTarballDecoder struct {
	reader     io.Reader
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
	waitGroup  sync.WaitGroup
}

func NewFirstTarballInTarballDecoder(reader io.Reader) *FirstTarballInTarballDecoder {
	pipeReader, pipeWriter := io.Pipe()
	return &FirstTarballInTarballDecoder{
		reader:     reader,
		pipeReader: pipeReader,
		pipeWriter: pipeWriter,
	}
}

func (self *FirstTarballInTarballDecoder) Decode() io.Reader {
	self.waitGroup.Add(1)
	go self.copyFirstTarball()
	return self.pipeReader
}

func (self *FirstTarballInTarballDecoder) Drain() {
	self.waitGroup.Wait()
}

func (self *FirstTarballInTarballDecoder) copyFirstTarball() {
	defer self.waitGroup.Done()

	if reader, err := OpenFirstTarballInTarball(self.reader); err == nil {
		if _, err := io.Copy(self.pipeWriter, reader); err == nil {
			self.pipeWriter.Close()
		} else {
			self.pipeWriter.CloseWithError(err)
		}
	} else {
		self.pipeWriter.CloseWithError(err)
	}
}
