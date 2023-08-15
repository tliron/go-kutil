package streampackage

import (
	"archive/tar"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zip"
)

func getRoot(t *testing.T) string {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("KUTIL_TEST_ROOT"); !ok {
		var err error
		if root, err = os.Getwd(); err != nil {
			t.Errorf("os.Getwd: %s", err.Error())
		}
	}
	return root
}

func readStreamPackage(t *testing.T, streamPackage StreamPackage) {
	defer func() {
		if err := streamPackage.Close(); err != nil {
			t.Errorf("streamPackage.Close: %s", err.Error())
		}
	}()

	for {
		if stream, err := streamPackage.Next(); err == nil {
			if stream == nil {
				break
			}

			if reader, path, _, err := stream.Open(context.TODO()); err == nil {
				t.Logf("reading stream: %s", path)

				if _, err := io.ReadAll(reader); err != nil {
					t.Errorf("io.ReadAll: %s", err.Error())
					return
				}

				if err := reader.Close(); err != nil {
					t.Errorf("reader.Close: %s", err.Error())
					return
				}
			} else {
				t.Errorf("stream.Open: %s", err.Error())
				return
			}
		} else {
			t.Errorf("streamPackage.Next: %s", err.Error())
			return
		}
	}
}

func createTarball(t *testing.T, tarWriter *tar.Writer) {
	defer func() {
		if err := tarWriter.Close(); err != nil {
			t.Errorf("tarWriter.Close: %s", err.Error())
		}
	}()

	sourcePath := filepath.Join(getRoot(t), "streampackage")

	if err := filepath.WalkDir(sourcePath, func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() {
			t.Logf("adding file to tarball: %s", path)

			if file, err := os.Open(path); err == nil {
				if stat, err := file.Stat(); err == nil {
					header := tar.Header{
						Name:    path[len(sourcePath):],
						Size:    stat.Size(),
						Mode:    int64(stat.Mode()),
						ModTime: stat.ModTime(),
					}

					if err := tarWriter.WriteHeader(&header); err == nil {
						if _, err := io.Copy(tarWriter, file); err != nil {
							t.Errorf("io.Copy: %s", err.Error())
							return nil
						}
					} else {
						t.Errorf("tarWriter.WriteHeader: %s", err.Error())
						return nil
					}
				} else {
					t.Errorf("file.Stat: %s", err.Error())
					return nil
				}
			} else {
				t.Errorf("os.Open: %s", err.Error())
				return nil
			}
		}
		return nil
	}); err != nil {
		t.Errorf("filepath.WalkDir: %s", err.Error())
		return
	}
}

func createZip(t *testing.T, zipWriter *zip.Writer) {
	sourcePath := filepath.Join(getRoot(t), "streampackage")

	if err := filepath.WalkDir(sourcePath, func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() {
			t.Logf("adding file to zip: %s", path)

			if file, err := os.Open(path); err == nil {
				if stat, err := file.Stat(); err == nil {
					header := zip.FileHeader{
						Name:     path[len(sourcePath):],
						Modified: stat.ModTime(),
					}

					if writer, err := zipWriter.CreateHeader(&header); err == nil {
						if _, err := io.Copy(writer, file); err != nil {
							t.Errorf("io.Copy: %s", err.Error())
							return nil
						}
					} else {
						t.Errorf("zipWriter.CreateHeader: %s", err.Error())
						return nil
					}
				} else {
					t.Errorf("file.Stat: %s", err.Error())
					return nil
				}
			} else {
				t.Errorf("os.Open: %s", err.Error())
				return nil
			}
		}
		return nil
	}); err != nil {
		t.Errorf("filepath.WalkDir: %s", err.Error())
	}
}
