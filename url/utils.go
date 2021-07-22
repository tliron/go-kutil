package url

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFormat(path string) string {
	extension := filepath.Ext(path)
	if extension == "" {
		return ""
	}
	extension = strings.ToLower(extension[1:])
	if extension == "yml" {
		extension = "yaml"
	}
	return extension
}

func DeleteTemporaryFile(path string) error {
	if err := os.Remove(path); err == nil {
		log.Infof("deleted temporary file %q", path)
		return nil
	} else if os.IsNotExist(err) {
		log.Infof("temporary file already deleted %q", path)
		return nil
	} else {
		log.Errorf("could not delete temporary file %q: %s", path, err.Error())
		return err
	}
}

func DeleteTemporaryDir(path string) error {
	if err := os.RemoveAll(path); err == nil {
		log.Infof("deleted temporary dir %q", path)
		return nil
	} else if os.IsNotExist(err) {
		log.Infof("temporary dir already deleted %q", path)
		return nil
	} else {
		log.Errorf("could not delete temporary dir %q: %s", path, err.Error())
		return err
	}
}
