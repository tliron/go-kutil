package streampackage

import (
	"io/fs"
	"path/filepath"
)

func NewDirStreamPackage(path string) (*StaticStreamPackage, error) {
	length := len(path)
	var streams []Stream
	if err := filepath.WalkDir(path, func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() {
			streams = append(streams, NewFileStream(path, path[length:]))
		}
		return nil
	}); err == nil {
		return NewStaticStreamPackage(streams...), nil
	} else {
		return nil, err
	}
}
