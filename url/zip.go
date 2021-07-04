package url

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"strings"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// ZipURL
//
// Inspired by Java's JarURLConnection:
// https://docs.oracle.com/javase/8/docs/api/java/net/JarURLConnection.html
//

type ZipURL struct {
	Path       string
	ArchiveURL URL
}

func NewZipURL(path string, archiveUrl URL) *ZipURL {
	// Must be absolute
	path = strings.TrimLeft(path, "/")

	return &ZipURL{
		Path:       path,
		ArchiveURL: archiveUrl,
	}
}

func NewValidZipURL(path string, archiveUrl URL) (*ZipURL, error) {
	self := NewZipURL(path, archiveUrl)
	if reader, err := self.OpenArchive(); err == nil {
		defer reader.Close()

		for _, file := range reader.Reader.File {
			if self.Path == file.Name {
				return self, nil
			}
		}

		return nil, fmt.Errorf("path %q not found in zip: %s", path, archiveUrl.String())
	} else {
		return nil, err
	}
}

func NewValidRelativeZipURL(path string, origin *ZipURL) (*ZipURL, error) {
	self := origin.Relative(path).(*ZipURL)
	if reader, err := self.OpenArchive(); err == nil {
		for _, file := range reader.Reader.File {
			if self.Path == file.Name {
				return self, nil
			}
		}

		return nil, fmt.Errorf("path %q not found in zip: %s", path, self.ArchiveURL.String())
	} else {
		return nil, err
	}
}

func ParseZipURL(url string, context *Context) (*ZipURL, error) {
	if archive, path, err := parseZipURL(url); err == nil {
		if archiveUrl, err := NewURL(archive, context); err == nil {
			return NewZipURL(path, archiveUrl), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ParseValidZipURL(url string, context *Context) (*ZipURL, error) {
	if archive, path, err := parseZipURL(url); err == nil {
		if zipUrl, err := NewURL(archive, context); err == nil {
			return NewValidZipURL(path, zipUrl)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// URL interface
// fmt.Stringer interface
func (self *ZipURL) String() string {
	return self.Key()
}

// URL interface
func (self *ZipURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *ZipURL) Origin() URL {
	path := pathpkg.Dir(self.Path)
	if path != "/" {
		path += "/"
	}

	// Note: deleteArchive is *not* copied over
	return &ZipURL{
		Path:       path,
		ArchiveURL: self.ArchiveURL,
	}
}

// URL interface
func (self *ZipURL) Relative(path string) URL {
	// Note: deleteArchive is *not* copied over
	return &ZipURL{
		Path:       pathpkg.Join(self.Path, path),
		ArchiveURL: self.ArchiveURL,
	}
}

// URL interface
func (self *ZipURL) Key() string {
	return fmt.Sprintf("zip:%s!/%s", self.ArchiveURL.String(), self.Path)
}

// URL interface
func (self *ZipURL) Open() (io.ReadCloser, error) {
	if archiveReader, err := self.OpenArchive(); err == nil {
		for _, file := range archiveReader.Reader.File {
			if self.Path == file.Name {
				if entryReader, err := file.Open(); err == nil {
					return NewZipEntryReader(entryReader, archiveReader), nil
				} else {
					archiveReader.Close()
					return nil, err
				}
			}
		}

		archiveReader.Close()
		return nil, fmt.Errorf("path %q not found in archive: %s", self.Path, self.ArchiveURL.String())
	} else {
		return nil, err
	}
}

// URL interface
func (self *ZipURL) Context() *Context {
	return self.ArchiveURL.Context()
}

func (self *ZipURL) OpenArchive() (*ZipReader, error) {
	if file, err := self.ArchiveURL.Context().Open(self.ArchiveURL); err == nil {
		return OpenZipFromFile(file)
	} else {
		return nil, err
	}
}

//
// ZipReader
//

type ZipReader struct {
	Reader *zip.Reader
	File   *os.File
}

func NewZipReader(reader *zip.Reader, file *os.File) *ZipReader {
	return &ZipReader{reader, file}
}

// io.Closer interface
func (self *ZipReader) Close() error {
	return self.File.Close()
}

//
// ZipEntryReader
//

type ZipEntryReader struct {
	EntryReader   io.ReadCloser
	ArchiveReader *ZipReader
}

func NewZipEntryReader(entryReader io.ReadCloser, archiveReader *ZipReader) *ZipEntryReader {
	return &ZipEntryReader{entryReader, archiveReader}
}

// io.Reader interface
func (self *ZipEntryReader) Read(p []byte) (n int, err error) {
	return self.EntryReader.Read(p)
}

// io.Closer interface
func (self *ZipEntryReader) Close() error {
	err1 := self.EntryReader.Close()
	err2 := self.ArchiveReader.Close()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

// Utils

func OpenZipFromFile(file *os.File) (*ZipReader, error) {
	if stat, err := file.Stat(); err == nil {
		size := stat.Size()
		if reader, err := zip.NewReader(file, size); err == nil {
			return NewZipReader(reader, file), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func parseZipURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "zip:") {
		if split := strings.Split(url[4:], "!"); len(split) == 2 {
			return split[0], split[1], nil
		} else {
			return "", "", fmt.Errorf("malformed \"zip:\" URL: %s", url)
		}
	} else {
		return "", "", fmt.Errorf("not a \"zip:\" URL: %s", url)
	}
}
