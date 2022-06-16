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
	files             map[string]string
	dirs              map[string]string
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

func (self *Context) OpenFile(url URL) (*os.File, error) {
	if path, err := self.GetLocalPath(url); err == nil {
		return os.Open(path)
	} else {
		return nil, err
	}
}

// Will download the file to the local temporary directory if not already locally available
func (self *Context) GetLocalPath(url URL) (string, error) {
	if fileUrl, ok := url.(*FileURL); ok {
		// No need to download file URLs
		return fileUrl.Path, nil
	}

	key := url.Key()

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.files != nil {
		if path, ok := self.files[key]; ok {
			if ok, err := util.DoesFileExist(path); err == nil {
				if ok {
					return path, nil
				} else {
					delete(self.files, key)
				}
			} else {
				return "", err
			}
		}
	}

	temporaryPathPattern := fmt.Sprintf("kutil-%s-*", util.SanitizeFilename(key))
	if file, err := Download(url, temporaryPathPattern); err == nil {
		if self.files == nil {
			self.files = make(map[string]string)
		}
		path := file.Name()
		self.files[key] = path
		return path, nil
	} else {
		return "", err
	}
}

func (self *Context) Release() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	var err error

	if self.files != nil {
		for _, path := range self.files {
			if err_ := DeleteTemporaryFile(path); err_ != nil {
				err = err_
			}
		}

		self.files = nil
	}

	if self.dirs != nil {
		for _, path := range self.dirs {
			if err_ := DeleteTemporaryDir(path); err_ != nil {
				err = err_
			}
		}

		self.files = nil
	}

	return err
}
