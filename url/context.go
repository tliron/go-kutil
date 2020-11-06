package url

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/tliron/kutil/util"
)

type Credentials struct {
	Username string
	Password string
	Token    string
}

//
// Context
//

type Context struct {
	paths             map[string]string
	httpRoundTrippers map[string]http.RoundTripper
	credentials       map[string]*Credentials
	lock              sync.Mutex // for paths
}

func NewContext() *Context {
	return &Context{}
}

// Not thread-safe
func (self *Context) SetHTTPRoundTripper(host string, httpRoundTripper http.RoundTripper) {
	if self.httpRoundTrippers == nil {
		self.httpRoundTrippers = make(map[string]http.RoundTripper)
	}
	self.httpRoundTrippers[host] = httpRoundTripper
}

// Not thread-safe
func (self *Context) GetHTTPRoundTripper(host string) http.RoundTripper {
	if self.httpRoundTrippers != nil {
		httpRoundTripper, _ := self.httpRoundTrippers[host]
		return httpRoundTripper
	} else {
		return nil
	}
}

// Not thread-safe
func (self *Context) SetCredentials(host string, username string, password string, token string) {
	if self.credentials == nil {
		self.credentials = make(map[string]*Credentials)
	}
	self.credentials[host] = &Credentials{
		Username: username,
		Password: password,
		Token:    token,
	}
}

// Not thread-safe
func (self *Context) GetCredentials(host string) *Credentials {
	if self.credentials != nil {
		credentials, _ := self.credentials[host]
		return credentials
	} else {
		return nil
	}
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

	temporaryPathPattern := fmt.Sprintf("puccini-%s-*", util.SanitizeFilename(key))
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
