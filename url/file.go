package url

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//
// FileURL
//

type FileURL struct {
	Path string

	context *Context
}

func NewFileURL(path string, context *Context) *FileURL {
	if context == nil {
		context = NewContext()
	}

	return &FileURL{
		Path:    path,
		context: context,
	}
}

func NewValidFileURL(path string, context *Context) (*FileURL, error) {
	isDir := strings.HasSuffix(path, "/")

	if filepath.IsAbs(path) {
		path = filepath.Clean(path)
	} else {
		var err error
		if path, err = filepath.Abs(path); err != nil {
			return nil, err
		}
	}

	if info, err := os.Stat(path); err == nil {
		if isDir {
			if !info.Mode().IsDir() {
				return nil, fmt.Errorf("URL path does not point to a directory: %s", path)
			}
		} else if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("URL path does not point to a file: %s", path)
		}
	} else {
		return nil, err
	}

	return NewFileURL(path, context), nil
}

func NewValidRelativeFileURL(path string, origin *FileURL) (*FileURL, error) {
	isDir := strings.HasSuffix(path, "/")
	path = filepath.Join(origin.Path, path)
	if isDir {
		path += "/"
	}
	return NewValidFileURL(path, origin.context)
}

// URL interface
// fmt.Stringer interface
func (self *FileURL) String() string {
	return self.Key()
}

// URL interface
func (self *FileURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *FileURL) Origin() URL {
	return &FileURL{
		Path:    filepath.Dir(self.Path),
		context: self.context,
	}
}

// URL interface
func (self *FileURL) Relative(path string) URL {
	return NewFileURL(filepath.Join(self.Path, path), self.context)
}

// URL interface
func (self *FileURL) Key() string {
	return "file:" + self.Path
}

// URL interface
func (self *FileURL) Open() (io.ReadCloser, error) {
	if reader, err := os.Open(self.Path); err == nil {
		return reader, nil
	} else {
		return nil, err
	}
}

// URL interface
func (self *FileURL) Context() *Context {
	return self.context
}

func isValidFile(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.Mode().IsRegular()
	} else {
		return false
	}
}
