package url

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"strings"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

// TODO: xz support, consider: https://github.com/ulikunitz/xz

var TARBALL_ARCHIVE_FORMATS = []string{"tar", "tar.gz"}

func IsValidTarballArchiveFormat(archiveFormat string) bool {
	for _, archiveFormat_ := range TARBALL_ARCHIVE_FORMATS {
		if archiveFormat_ == archiveFormat {
			return true
		}
	}
	return false
}

//
// TarballURL
//
// Inspired by Java's JarURLConnection:
// https://docs.oracle.com/javase/8/docs/api/java/net/JarURLConnection.html
//

type TarballURL struct {
	Path          string
	ArchiveURL    URL
	ArchiveFormat string
}

func NewTarballURL(path string, archiveUrl URL, archiveFormat string) *TarballURL {
	path = strings.TrimLeft(path, "/")

	if archiveFormat == "" {
		archiveFormat = archiveUrl.Format()
	}

	return &TarballURL{
		Path:          path,
		ArchiveURL:    archiveUrl,
		ArchiveFormat: archiveFormat,
	}
}

func NewValidTarballURL(path string, archiveUrl URL, archiveFormat string) (*TarballURL, error) {
	self := NewTarballURL(path, archiveUrl, archiveFormat)
	if tarballReader, err := self.OpenArchive(); err == nil {
		defer tarballReader.Close()

		for {
			if header, err := tarballReader.TarReader.Next(); err == nil {
				if self.Path == fixTarballEntryPath(header.Name) {
					return self, nil
				}
			} else if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		return nil, NewNotFoundf("path %q not found in tarball: %s", path, archiveUrl.String())
	} else {
		return nil, err
	}
}

func NewValidRelativeTarballURL(path string, origin *TarballURL) (*TarballURL, error) {
	self := origin.Relative(path).(*TarballURL)
	if tarballReader, err := self.OpenArchive(); err == nil {
		defer tarballReader.Close()

		for {
			if header, err := tarballReader.TarReader.Next(); err == nil {
				if self.Path == fixTarballEntryPath(header.Name) {
					return self, nil
				}
			} else if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		return nil, NewNotFoundf("path %q not found in tarball: %s", self.Path, self.ArchiveURL.String())
	} else {
		return nil, err
	}
}

func ParseTarballURL(url string, context *Context) (*TarballURL, error) {
	if archiveUrl, path, err := parseTarballURL(url); err == nil {
		if archiveUrl_, err := NewURL(archiveUrl, context); err == nil {
			return NewTarballURL(path, archiveUrl_, ""), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ParseValidTarballURL(url string, context *Context) (*TarballURL, error) {
	if archiveUrl, path, err := parseTarballURL(url); err == nil {
		if archiveUrl_, err := NewURL(archiveUrl, context); err == nil {
			return NewValidTarballURL(path, archiveUrl_, "")
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// URL interface
// fmt.Stringer interface
func (self *TarballURL) String() string {
	return self.Key()
}

// URL interface
func (self *TarballURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *TarballURL) Origin() URL {
	path := pathpkg.Dir(self.Path)
	if path != "/" {
		path += "/"
	}

	return &TarballURL{
		Path:          path,
		ArchiveURL:    self.ArchiveURL,
		ArchiveFormat: self.ArchiveFormat,
	}
}

// URL interface
func (self *TarballURL) Relative(path string) URL {
	return &TarballURL{
		Path:          pathpkg.Join(self.Path, path),
		ArchiveURL:    self.ArchiveURL,
		ArchiveFormat: self.ArchiveFormat,
	}
}

// URL interface
func (self *TarballURL) Key() string {
	return fmt.Sprintf("tar:%s!/%s", self.ArchiveURL.String(), self.Path)
}

// URL interface
func (self *TarballURL) Open() (io.ReadCloser, error) {
	if tarballReader, err := self.OpenArchive(); err == nil {
		if tarballEntryReader, err := tarballReader.Open(self.Path); err == nil {
			if tarballEntryReader != nil {
				return tarballEntryReader, nil
			} else {
				tarballReader.Close()
				return nil, NewNotFoundf("path %q not found in archive: %s", self.Path, self.ArchiveURL.String())
			}
		} else {
			tarballReader.Close()
			return nil, err
		}
	} else {
		return nil, err
	}
}

// URL interface
func (self *TarballURL) Context() *Context {
	return self.ArchiveURL.Context()
}

func (self *TarballURL) OpenArchive() (*TarballReader, error) {
	if !IsValidTarballArchiveFormat(self.ArchiveFormat) {
		return nil, fmt.Errorf("unsupported tarball archive format: %q", self.ArchiveFormat)
	}

	if archiveReader, err := self.ArchiveURL.Open(); err == nil {
		switch self.ArchiveFormat {
		case "tar.gz":
			if gzipReader, err := gzip.NewReader(archiveReader); err == nil {
				return NewTarballReader(tar.NewReader(gzipReader), archiveReader, gzipReader), nil
			} else {
				archiveReader.Close()
				return nil, err
			}

		default:
			return NewTarballReader(tar.NewReader(archiveReader), archiveReader, nil), nil
		}
	} else {
		return nil, err
	}
}

//
// TarballReader
//

type TarballReader struct {
	TarReader         *tar.Reader
	ArchiveReader     io.ReadCloser
	CompressionReader io.ReadCloser
}

func NewTarballReader(reader *tar.Reader, archiveReader io.ReadCloser, compressionReader io.ReadCloser) *TarballReader {
	return &TarballReader{reader, archiveReader, compressionReader}
}

// io.Closer interface
func (self *TarballReader) Close() error {
	var err1 error
	if self.CompressionReader != nil {
		err1 = self.CompressionReader.Close()
	}
	err2 := self.ArchiveReader.Close()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

func (self *TarballReader) Open(path string) (*TarballEntryReader, error) {
	for {
		if header, err := self.TarReader.Next(); err == nil {
			if path == fixTarballEntryPath(header.Name) {
				return NewTarballEntryReader(self), nil
			}
		} else if err == io.EOF {
			break
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (self *TarballReader) Has(path string) (bool, error) {
	for {
		if header, err := self.TarReader.Next(); err == nil {
			if path == fixTarballEntryPath(header.Name) {
				return true, nil
			}
		} else if err == io.EOF {
			break
		} else {
			return false, err
		}
	}
	return false, nil
}

func (self *TarballReader) Iterate(f func(*tar.Header) bool) error {
	for {
		if header, err := self.TarReader.Next(); err == nil {
			if !f(header) {
				return nil
			}
		} else if err == io.EOF {
			break
		} else {
			return err
		}
	}
	return nil
}

//
// TarballEntryReader
//

type TarballEntryReader struct {
	TarballReader *TarballReader
}

func NewTarballEntryReader(tarballReader *TarballReader) *TarballEntryReader {
	return &TarballEntryReader{tarballReader}
}

// io.Reader interface
func (self *TarballEntryReader) Read(p []byte) (n int, err error) {
	return self.TarballReader.TarReader.Read(p)
}

// io.Closer interface
func (self *TarballEntryReader) Close() error {
	return self.TarballReader.Close()
}

// Utils

func OpenTarballFromFile(file *os.File) (*TarballReader, error) {
	return NewTarballReader(tar.NewReader(file), file, nil), nil // BAD
}

func parseTarballURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "tar:") {
		if split := strings.Split(url[4:], "!"); len(split) == 2 {
			return split[0], split[1], nil
		} else {
			return "", "", fmt.Errorf("malformed \"tar:\" URL: %s", url)
		}
	} else {
		return "", "", fmt.Errorf("not a \"tar:\" URL: %s", url)
	}
}

func fixTarballEntryPath(path string) string {
	if strings.HasPrefix(path, "./") {
		return path[3:]
	}
	return path
}
