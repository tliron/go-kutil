package streampackage

import (
	"context"
	"os"
	"testing"

	"github.com/klauspost/compress/zip"
	"github.com/tliron/exturl"
)

func TestZip(t *testing.T) {
	zipFile, err := os.CreateTemp("", "kutil-*")
	if err != nil {
		t.Errorf("os.CreateTemp: %s", err.Error())
	}
	zipPath := zipFile.Name()

	zipWriter := zip.NewWriter(zipFile)

	urlContext := exturl.NewContext()

	defer func() {
		if zipWriter != nil {
			if err := zipWriter.Close(); err != nil {
				t.Errorf("zipWriter.Close: %s", err.Error())
			}
		}
		if zipFile != nil {
			if err := zipFile.Close(); err != nil {
				t.Errorf("zipFile.Close: %s", err.Error())
			}
		}
		if err := os.Remove(zipPath); err != nil {
			t.Errorf("os.Remove: %s", err.Error())
		}
		if err := urlContext.Release(); err != nil {
			t.Errorf("urlContext.Release: %s", err.Error())
		}
	}()

	createZip(t, zipWriter)

	if err := zipWriter.Close(); err != nil {
		zipWriter = nil
		t.Errorf("zipWriter.Close: %s", err.Error())
		return
	}
	zipWriter = nil

	if err := zipFile.Close(); err != nil {
		zipFile = nil
		t.Errorf("zipFile.Close: %s", err.Error())
		return
	}
	zipFile = nil

	url := urlContext.NewFileURL(zipPath)
	if streamPackage, err := NewStreamPackage(context.TODO(), url, "zip"); err == nil {
		readStreamPackage(t, streamPackage)
	} else {
		t.Errorf("NewStreamPackage: %s", err.Error())
		return
	}
}
