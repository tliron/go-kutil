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
func (self *URLStream) Open(context contextpkg.Context) (io.ReadCloser, string, bool, error) {
	var err error
	if self.reader, err = self.url.Open(context); err == nil {
		if path, err := exturl.GetPath(self.url); err == nil {
			return self.reader, filepath.Base(path), false, nil
		} else {
			return nil, "", false, err
		}
	} else {
		return nil, "", false, err
	}
}
