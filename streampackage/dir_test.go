package streampackage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/tliron/exturl"
)

func TestDir(t *testing.T) {
	urlContext := exturl.NewContext()

	defer func() {
		if err := urlContext.Release(); err != nil {
			t.Logf("urlContext.Release: %s", err.Error())
		}
	}()

	url := urlContext.NewFileURL(filepath.Join(getRoot(t), "streampackage"))
	if streamPackage, err := NewStreamPackage(context.TODO(), url, ""); err == nil {
		readStreamPackage(t, streamPackage)
	} else {
		t.Errorf("NewStreamPackage: %s", err.Error())
	}
}
