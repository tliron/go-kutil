package url

import (
	"bytes"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"sync"

	"github.com/segmentio/ksuid"
	"github.com/tliron/kutil/util"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

var internal sync.Map

// `content` must be []byte or string
func RegisterInternalURL(path string, content any) error {
	if _, loaded := internal.LoadOrStore(path, util.ToBytes(content)); !loaded {
		return nil
	} else {
		return fmt.Errorf("internal URL conflict: %s", path)
	}
}

func DeregisterInternalURL(path string) {
	internal.Delete(path)
}

func UpdateInternalURL(path string, content any) {
	internal.Store(path, util.ToBytes(content))
}

func ReadToInternalURL(path string, reader io.Reader, context *Context) (*InternalURL, error) {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	if buffer, err := io.ReadAll(reader); err == nil {
		if err = RegisterInternalURL(path, buffer); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return NewValidInternalURL(path, context)
}

func ReadToInternalURLFromStdin(format string, context *Context) (*InternalURL, error) {
	path := fmt.Sprintf("<stdin:%s>", ksuid.New().String())
	if format != "" {
		path = fmt.Sprintf("%s.%s", path, format)
	}
	return ReadToInternalURL(path, os.Stdin, context)
}

//
// InternalURL
//

type InternalURL struct {
	Path    string
	Content []byte

	context *Context
}

func NewInternalURL(path string, context *Context) *InternalURL {
	if context == nil {
		context = NewContext()
	}

	return &InternalURL{
		Path:    path,
		context: context,
	}
}

func NewValidInternalURL(path string, context *Context) (*InternalURL, error) {
	if content, ok := internal.Load(path); ok {
		if context == nil {
			context = NewContext()
		}

		return &InternalURL{
			Path:    path,
			Content: content.([]byte),
			context: context,
		}, nil
	} else {
		return nil, fmt.Errorf("internal URL not found: %s", path)
	}
}

func NewValidRelativeInternalURL(path string, origin *InternalURL) (*InternalURL, error) {
	return NewValidInternalURL(pathpkg.Join(origin.Path, path), origin.context)
}

func (self *InternalURL) SetContent(content any) {
	self.Content = util.ToBytes(content)
}

// URL interface
// fmt.Stringer interface
func (self *InternalURL) String() string {
	return self.Key()
}

// URL interface
func (self *InternalURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *InternalURL) Origin() URL {
	path := pathpkg.Dir(self.Path)
	if path != "/" {
		path += "/"
	}

	return &InternalURL{
		Path:    path,
		context: self.context,
	}
}

// URL interface
func (self *InternalURL) Relative(path string) URL {
	return NewInternalURL(pathpkg.Join(self.Path, path), self.context)
}

// URL interface
func (self *InternalURL) Key() string {
	return "internal:" + self.Path
}

// URL interface
func (self *InternalURL) Open() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(self.Content)), nil
}

// URL interface
func (self *InternalURL) Context() *Context {
	return self.context
}
