package streampackage

import (
	contextpkg "context"
	"io"
	"path/filepath"

	"github.com/tliron/exturl"
)

//
// URLStream
//

type URLStream struct {
	url exturl.URL

	reader io.ReadCloser
}

func NewURLStream(url exturl.URL) *URLStream {
	return &URLStream{
		url: url,
	}
}

// Stream interface
func (self *URLStream) Open(context contextpkg.Context) (string, bool, io.Reader, error) {
	var err error
	if self.reader, err = self.url.Open(context); err == nil {
		if path, err := exturl.GetPath(self.url); err == nil {
			return filepath.Base(path), false, self.reader, nil
		} else {
			return "", false, nil, err
		}
	} else {
		return "", false, nil, err
	}
}

// Stream interface
func (self *URLStream) Close() error {
	return self.reader.Close()
}
