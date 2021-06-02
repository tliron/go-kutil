package fswatch

import (
	"errors"

	urlpkg "github.com/tliron/kutil/url"
)

type Watcher struct{}

func NewWatcher(context *urlpkg.Context) (*Watcher, error) {
	return nil, errors.New("watching is not supported on this platform")
}

func (self *Watcher) Add(path string) error {
	return nil
}

func (self *Watcher) Close() error {
	return nil
}

func (self *Watcher) Start(onChanged OnChangedFunc) {
}
