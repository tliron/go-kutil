package url

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/tliron/kutil/util"
)

//
// Context
//

type Context struct {
	paths            map[string]string
	httpRoundTripper http.RoundTripper
	lock             sync.Mutex
}

func NewContext() *Context {
	return &Context{}
}

// Not thread-safe
func (self *Context) SetHTTPRoundTripper(httpRoundTripper http.RoundTripper) {
	self.httpRoundTripper = httpRoundTripper
}

// Not thread-safe
func (self *Context) GetHTTPRoundTripper() http.RoundTripper {
	return self.httpRoundTripper
}

func (self *Context) Open(url URL) (*os.File, error) {
	if fileUrl, ok := url.(*FileURL); ok {
		// No need to download file URLs
		return os.Open(fileUrl.Path)
	}

	key := url.Key()

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.paths != nil {
		if path, ok := self.paths[key]; ok {
			if file, err := os.Open(path); err == nil {
				return file, nil
			} else if os.IsNotExist(err) {
				delete(self.paths, key)
			} else {
				return nil, err
			}
		}
	}

	// TODO: remove .zip?
	temporaryPathPattern := fmt.Sprintf("puccini-%s-*.zip", util.SanitizeFilename(key))
	if file, err := Download(url, temporaryPathPattern); err == nil {
		if self.paths == nil {
			self.paths = make(map[string]string)
		}
		self.paths[key] = file.Name()
		return file, nil
	} else {
		return nil, err
	}
}

func (self *Context) Release() {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.paths != nil {
		for _, path := range self.paths {
			DeleteTemporaryFile(path)
		}

		self.paths = nil
	}
}
