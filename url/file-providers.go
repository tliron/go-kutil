package url

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tliron/kutil/util"
)

//
// FileProvider
//

type FileProvider interface {
	Open() (string, bool, io.Reader, error)
	Close() error
}

//
// FileProviders
//

type FileProviders interface {
	Next() (FileProvider, error)
	Close() error
}

// `unpack` can be "tgz" or "zip"
func NewFileProviders(url URL, unpack string) (FileProviders, error) {
	var isFile bool
	var isDir bool

	if fileUrl, ok := url.(*FileURL); ok {
		isFile = true
		if stat, err := os.Stat(fileUrl.Path); err == nil {
			isDir = stat.IsDir()
		} else {
			return nil, err
		}
	}

	if isDir {
		return NewDirFileProviders(url.(*FileURL).Path)
	} else {
		switch unpack {
		case "tar":
			return NewTarFileProviders(url)

		case "tgz":
			return NewTarGZipFileProviders(url)

		case "zip":
			if path, err := url.Context().GetLocalPath(url); err == nil {
				return NewZipFileProviders(path)
			} else {
				return nil, err
			}
		}

		if isFile {
			path := url.(*FileURL).Path
			return NewStaticFileProviders(NewFileFileProvider(path, filepath.Base(path))), nil
		} else {
			return NewStaticFileProviders(NewURLFileProvider(url)), nil
		}
	}
}

//
// StaticFileProviders
//

type StaticFileProviders struct {
	sources []FileProvider

	index int
}

func NewStaticFileProviders(sources ...FileProvider) *StaticFileProviders {
	return &StaticFileProviders{
		sources: sources,
	}
}

// FileProviders interface
func (self *StaticFileProviders) Next() (FileProvider, error) {
	if self.index < len(self.sources) {
		source := self.sources[self.index]
		self.index++
		return source, nil
	} else {
		return nil, nil
	}
}

// FileProviders interface
func (self *StaticFileProviders) Close() error {
	return nil
}

//
// FileFileProvider
//

type FileFileProvider struct {
	localPath    string
	providedPath string

	file *os.File
}

func NewFileFileProvider(localPath string, providedPath string) *FileFileProvider {
	return &FileFileProvider{
		localPath:    localPath,
		providedPath: providedPath,
	}
}

func NewDirFileProviders(path string) (FileProviders, error) {
	length := len(path)
	var sources []FileProvider
	_ = filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			sources = append(sources, NewFileFileProvider(path, path[length:]))
		}
		return nil
	})
	return NewStaticFileProviders(sources...), nil
}

// FileReader interface
func (self *FileFileProvider) Open() (string, bool, io.Reader, error) {
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

// FileReader interface
func (self *FileFileProvider) Close() error {
	return self.file.Close()
}

//
// URLFileProvider
//

type URLFileProvider struct {
	url URL

	reader io.ReadCloser
}

func NewURLFileProvider(url URL) *URLFileProvider {
	return &URLFileProvider{
		url: url,
	}
}

// FileReader interface
func (self *URLFileProvider) Open() (string, bool, io.Reader, error) {
	var err error
	if self.reader, err = self.url.Open(); err == nil {
		if path, err := GetPath(self.url); err == nil {
			return filepath.Base(path), false, self.reader, nil
		} else {
			return "", false, nil, err
		}
	} else {
		return "", false, nil, err
	}
}

// FileReader interface
func (self *URLFileProvider) Close() error {
	return self.reader.Close()
}

//
// TarFileProvider
//

type TarFileProvider struct {
	header    *tar.Header
	tarReader *tar.Reader
}

func NewTarFileReader(header *tar.Header, tarReader *tar.Reader) *TarFileProvider {
	return &TarFileProvider{
		header:    header,
		tarReader: tarReader,
	}
}

// FileReader interface
func (self *TarFileProvider) Open() (string, bool, io.Reader, error) {
	return self.header.Name, util.IsFileExecutable(fs.FileMode(self.header.Mode)), self.tarReader, nil
}

// FileReader interface
func (self *TarFileProvider) Close() error {
	return nil
}

//
// TarFileProviders
//

type TarFileProviders struct {
	reader    io.ReadCloser
	tarReader *tar.Reader
}

func NewTarFileProviders(url URL) (*TarFileProviders, error) {
	var self TarFileProviders

	var err error
	if self.reader, err = url.Open(); err == nil {
		self.tarReader = tar.NewReader(self.reader)
	} else {
		return nil, err
	}

	return &self, nil
}

// FileProviders interface
func (self *TarFileProviders) Next() (FileProvider, error) {
	for {
		header, err := self.tarReader.Next()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			} else {
				return nil, err
			}
		}
		if header.Typeflag == tar.TypeReg {
			return NewTarFileReader(header, self.tarReader), nil
		}
	}
}

// FileProviders interface
func (self *TarFileProviders) Close() error {
	return self.reader.Close()
}

//
// TarGZipFileProviders
//

type TarGZipFileProviders struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
	tarReader  *tar.Reader
}

func NewTarGZipFileProviders(url URL) (*TarGZipFileProviders, error) {
	var self TarGZipFileProviders

	var err error
	if self.reader, err = url.Open(); err == nil {
		if self.gzipReader, err = gzip.NewReader(self.reader); err == nil {
			self.tarReader = tar.NewReader(self.gzipReader)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &self, nil
}

// FileProviders interface
func (self *TarGZipFileProviders) Next() (FileProvider, error) {
	for {
		header, err := self.tarReader.Next()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			} else {
				return nil, err
			}
		}
		if header.Typeflag == tar.TypeReg {
			return NewTarFileReader(header, self.tarReader), nil
		}
	}
}

// FileProviders interface
func (self *TarGZipFileProviders) Close() error {
	self.gzipReader.Close() // TODO: err?
	return self.reader.Close()
}

//
// ZipFileProvider
//

type ZipFileProvider struct {
	file *zip.File

	reader io.ReadCloser
}

func NewZipFileProvider(file *zip.File) *ZipFileProvider {
	return &ZipFileProvider{
		file: file,
	}
}

// FileReader interface
func (self *ZipFileProvider) Open() (string, bool, io.Reader, error) {
	var err error
	if self.reader, err = self.file.Open(); err == nil {
		return self.file.Name, util.IsFileExecutable(self.file.Mode()), self.reader, nil
	} else {
		return "", false, nil, err
	}
}

// FileReader interface
func (self *ZipFileProvider) Close() error {
	return self.reader.Close()
}

//
// ZipFileProviders
//

type ZipFileProviders struct {
	zipReader *zip.ReadCloser

	sources []*ZipFileProvider
	index   int
}

func NewZipFileProviders(path string) (*ZipFileProviders, error) {
	var self ZipFileProviders
	var err error
	if self.zipReader, err = zip.OpenReader(path); err == nil {
		for _, file := range self.zipReader.File {
			if !file.FileInfo().IsDir() {
				self.sources = append(self.sources, NewZipFileProvider(file))
			}
		}
		return &self, nil
	} else {
		return nil, err
	}
}

// FileProviders interface
func (self *ZipFileProviders) Next() (FileProvider, error) {
	if self.index < len(self.sources) {
		source := self.sources[self.index]
		self.index++
		return source, nil
	} else {
		return nil, nil
	}
}

// FileProviders interface
func (self *ZipFileProviders) Close() error {
	return self.zipReader.Close()
}
