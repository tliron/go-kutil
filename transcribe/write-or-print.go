package transcribe

import (
	"io"
	"os"
	"path/filepath"

	"github.com/tliron/kutil/terminal"
)

const DIRECTORY_WRITE_PERMISSIONS = 0700
const FILE_WRITE_PERMISSIONS = 0600

func WriteOrPrint(value any, format string, writer io.Writer, strict bool, pretty bool, output string) error {
	if output != "" {
		if f, err := OpenFileForWrite(output); err == nil {
			defer f.Close()
			return Write(value, format, terminal.Indent, strict, f)
		} else {
			return err
		}
	} else {
		return Print(value, format, writer, strict, pretty)
	}
}

func OpenFileForWrite(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), DIRECTORY_WRITE_PERMISSIONS); err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FILE_WRITE_PERMISSIONS)
	} else {
		return nil, err
	}
}
