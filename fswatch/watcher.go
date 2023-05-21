//go:build !wasm

package fswatch

import (
	"github.com/fsnotify/fsnotify"
	"github.com/tliron/exturl"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	context *exturl.Context
}

func NewWatcher(context *exturl.Context) (*Watcher, error) {
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		return &Watcher{
			watcher: watcher,
			context: context,
		}, nil
	} else {
		return nil, err
	}
}

func (self *Watcher) Add(path string) error {
	return self.watcher.Add(path)
}

func (self *Watcher) Close() error {
	return self.watcher.Close()
}

func (self *Watcher) Start(onChanged OnChangedFunc) {
	go func() {
		for {
			select {
			case event, ok := <-self.watcher.Events:
				if !ok {
					return
				}

				onChanged(exturl.NewFileURL(event.Name, self.context))

			case err, ok := <-self.watcher.Errors:
				if !ok {
					return
				}

				log.Errorf("%s", err.Error())
			}
		}
	}()
}
