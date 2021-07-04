package url

import (
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"strings"
	"sync"

	"github.com/segmentio/ksuid"
	"github.com/tliron/kutil/util"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

var internal sync.Map

func RegisterInternalURL(path string, content string) error {
	if _, loaded := internal.LoadOrStore(path, content); !loaded {
		return nil
	} else {
		return fmt.Errorf("internal URL conflict: %s", path)
	}
}

func DeregisterInternalURL(path string) {
	internal.Delete(path)
}

func UpdateInternalURL(path string, content string) {
	internal.Store(path, content)
}

func ReadToInternalURL(path string, reader io.Reader) (*InternalURL, error) {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	if buffer, err := io.ReadAll(reader); err == nil {
		if err = RegisterInternalURL(path, util.BytesToString(buffer)); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return NewValidInternalURL(path, nil)
}

func ReadToInternalURLFromStdin(format string) (*InternalURL, error) {
	path := fmt.Sprintf("<stdin:%s>", ksuid.New().String())
	if format != "" {
		path = fmt.Sprintf("%s.%s", path, format)
	}
	return ReadToInternalURL(path, os.Stdin)
}

//
// InternalURL
//

type InternalURL struct {
	Path    string
	Content string

	context *Context
}

func NewInternalURL(path string, context *Context) *InternalURL {
	if context == nil {
		context = NewContext()
	}

	return &InternalURL{
		Path:    path,
		Content: "",
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
			Content: content.(string),
			context: context,
		}, nil
	} else {
		return nil, fmt.Errorf("internal URL not found: %s", path)
	}
}

func NewValidRelativeInternalURL(path string, origin *InternalURL) (*InternalURL, error) {
	return NewValidInternalURL(pathpkg.Join(origin.Path, path), origin.context)
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
		Content: "",
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
	return io.NopCloser(strings.NewReader(self.Content)), nil
}

// URL interface
func (self *InternalURL) Context() *Context {
	return self.context
}
