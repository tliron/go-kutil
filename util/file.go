package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DoesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func IsFileHidden(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func IsFileExecutable(mode fs.FileMode) bool {
	return mode&0100 != 0
}
