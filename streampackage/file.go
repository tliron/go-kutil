package streampackage

import (
	contextpkg "context"
	"io"
	"os"

	"github.com/tliron/kutil/util"
)

//
// FileStream
//

type FileStream struct {
	localPath string
	path      string

	file *os.File
}

func NewFileStream(localPath string, path string) *FileStream {
	return &FileStream{
		localPath: localPath,
		path:      path,
	}
}

// Stream interface
func (self *FileStream) Open(context contextpkg.Context) (io.ReadCloser, string, bool, error) {
	if stat, err := os.Stat(self.localPath); err == nil {
		if self.file, err = os.Open(self.localPath); err == nil {
			return self.file, self.path, util.IsFileExecutable(stat.Mode()), nil
		} else {
			return nil, "", false, err
		}
	} else {
		return nil, "", false, err
	}
}
