package url

import (
	"io"
	"os"

	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
)

var log = logging.MustGetLogger("puccini.url")

func ReadString(url URL) (string, error) {
	if reader, err := url.Open(); err == nil {
		defer reader.Close()
		buffer, err := io.ReadAll(reader)
		return util.BytesToString(buffer), err
	} else {
		return "", err
	}
}

func ReadARD(url URL, locate bool) (ard.Value, ard.Locator, error) {
	if reader, err := url.Open(); err == nil {
		defer reader.Close()
		return ard.Read(reader, url.Format(), locate)
	} else {
		return nil, nil, err
	}
}

func Size(url URL) (int64, error) {
	if reader, err := url.Open(); err == nil {
		defer reader.Close()
		return util.ReaderSize(reader)
	} else {
		return 0, err
	}
}

func DownloadTo(url URL, path string) error {
	if writer, err := os.Create(path); err == nil {
		if reader, err := url.Open(); err == nil {
			defer reader.Close()
			log.Infof("downloading from %q to file %q", url.String(), path)
			if _, err = io.Copy(writer, reader); err == nil {
				return nil
			} else {
				log.Warningf("failed to download from %q", url.String())
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func Download(url URL, temporaryPathPattern string) (*os.File, error) {
	if file, err := os.CreateTemp("", temporaryPathPattern); err == nil {
		path := file.Name()
		if reader, err := url.Open(); err == nil {
			defer reader.Close()
			log.Infof("downloading from %q to temporary file %q", url.String(), path)
			if _, err = io.Copy(file, reader); err == nil {
				atexit.Register(func() {
					DeleteTemporaryFile(path)
				})
				return file, nil
			} else {
				log.Warningf("failed to download from %q", url.String())
				DeleteTemporaryFile(path)
				return nil, err
			}
		} else {
			DeleteTemporaryFile(path)
			return nil, err
		}
	} else {
		return nil, err
	}
}
