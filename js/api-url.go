package js

import (
	"fmt"

	urlpkg "github.com/tliron/kutil/url"
)

type URLAPI struct {
	Url urlpkg.URL
}

func (self URLAPI) Directory() (string, error) {
	origin := self.Url.Origin()
	if origin_, ok := origin.(*urlpkg.FileURL); ok {
		return origin_.Path, nil
	} else {
		return "", fmt.Errorf("not a file: %s", origin)
	}
}

func (self URLAPI) GetRelativeURL(url string) (urlpkg.URL, error) {
	var origins []urlpkg.URL
	var context *urlpkg.Context
	if self.Url != nil {
		origins = []urlpkg.URL{self.Url.Origin()}
		context = self.Url.Context()
	}

	return urlpkg.NewValidURL(url, origins, context)
}

func (self URLAPI) Load(url string) (string, error) {
	if url_, err := self.GetRelativeURL(url); err == nil {
		return urlpkg.ReadString(url_)
	} else {
		return "", err
	}
}

func (self URLAPI) Download(sourceUrl string, targetPath string) error {
	if sourceUrl_, err := self.GetRelativeURL(sourceUrl); err == nil {
		return urlpkg.DownloadTo(sourceUrl_, targetPath)
	} else {
		return err
	}
}
