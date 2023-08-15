package streampackage

import (
	"archive/tar"
	"context"
	"os"
	"testing"

	"github.com/tliron/exturl"
)

func TestTar(t *testing.T) {
	tarballFile, err := os.CreateTemp("", "kutil-*")
	if err != nil {
		t.Errorf("os.CreateTemp: %s", err.Error())
	}
	tarballPath := tarballFile.Name()

	urlContext := exturl.NewContext()

	defer func() {
		if tarballFile != nil {
			if err := tarballFile.Close(); err != nil {
				t.Errorf("tarFile.Close: %s", err.Error())
			}
		}
		if err := os.Remove(tarballPath); err != nil {
			t.Errorf("os.Remove: %s", err.Error())
		}
		if err := urlContext.Release(); err != nil {
			t.Errorf("urlContext.Release: %s", err.Error())
		}
	}()

	createTarball(t, tar.NewWriter(tarballFile))

	if err := tarballFile.Close(); err != nil {
		tarballFile = nil
		t.Errorf("tarFile.Close: %s", err.Error())
		return
	}
	tarballFile = nil

	url := urlContext.NewFileURL(tarballPath)
	if streamPackage, err := NewStreamPackage(context.TODO(), url, "tar"); err == nil {
		readStreamPackage(t, streamPackage)
	} else {
		t.Errorf("NewStreamPackage: %s", err.Error())
		return
	}
}
