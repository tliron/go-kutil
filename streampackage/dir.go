package streampackage

import (
	contextpkg "context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tliron/kutil/util"
)

//
// DirStream
//

type DirStream struct {
	localPath    string
	providedPath string

	file *os.File
}

func NewDirStream(localPath string, providedPath string) *DirStream {
	return &DirStream{
		localPath:    localPath,
		providedPath: providedPath,
	}
}

func NewDirStreamPackage(path string) (StreamPackage, error) {
	length := len(path)
	var sources []Stream
	if err := filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			sources = append(sources, NewDirStream(path, path[length:]))
		}
		return nil
	}); err == nil {
		return NewStaticStreamPackage(sources...), nil
	} else {
		return nil, err
	}
}

// Stream interface
func (self *DirStream) Open(context contextpkg.Context) (string, bool, io.Reader, error) {
	if stat, err := os.Stat(self.localPath); err == nil {
		if self.file, err = os.Open(self.localPath); err == nil {
			return self.providedPath, util.IsFileExecutable(stat.Mode()), self.file, nil
		} else {
			return "", false, nil, err
		}
	} else {
		return "", false, nil, err
	}
}

// Stream interface
func (self *DirStream) Close() error {
	return self.file.Close()
}
