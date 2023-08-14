package streampackage

import (
	contextpkg "context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tliron/exturl"
)

// `archiveFormat` can be "tar", "tgz", or "zip"
func NewStreamPackage(context contextpkg.Context, url exturl.URL, archiveFormat string) (StreamPackage, error) {
	var isFile bool
	var isDir bool

	if fileUrl, ok := url.(*exturl.FileURL); ok {
		isFile = true
		if stat, err := os.Stat(fileUrl.Path); err == nil {
			isDir = stat.IsDir()
		} else {
			return nil, err
		}
	}

	if isDir {
		return NewDirStreamPackage(url.(*exturl.FileURL).Path)
	} else {
		switch archiveFormat {
		case "":
			if isFile {
				path := url.(*exturl.FileURL).Path
				return NewStaticStreamPackage(NewFileStream(path, filepath.Base(path))), nil
			} else {
				return NewStaticStreamPackage(NewURLStream(url)), nil
			}

		case "tar":
			return NewTarStreamPackage(context, url)

		case "tgz":
			return NewTarGZipStreamPackage(context, url)

		case "zip":
			if path, err := url.Context().GetLocalPath(context, url); err == nil {
				return NewZipStreamPackage(path)
			} else {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unsupported archive format: %s", archiveFormat)
		}
	}
}
