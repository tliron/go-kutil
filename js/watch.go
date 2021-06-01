// +build !wasm

package js

import (
	"github.com/fsnotify/fsnotify"
	urlpkg "github.com/tliron/kutil/url"
)

type Watcher struct {
	watcher *fsnotify.Watcher
}

func NewWatcher() (*Watcher, error) {
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		return &Watcher{watcher: watcher}, nil
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

func (self *Environment) Watch(onChanged OnChangedFunc) error {
	var err error
	if self.Watcher, err = NewWatcher(); err == nil {
		go func() {
			for {
				select {
				case event, ok := <-self.Watcher.watcher.Events:
					if !ok {
						return
					}

					id := urlpkg.NewFileURL(event.Name, nil).Key()
					var module *Module
					if module_ := self.Modules.Get(id); module_ != nil {
						module = module_.Export().(*Module)
					}
					onChanged(id, module)

				case err, ok := <-self.Watcher.watcher.Errors:
					if !ok {
						return
					}

					self.Log.Errorf("%s", err.Error())
				}
			}
		}()

		return nil
	} else {
		return err
	}
}
