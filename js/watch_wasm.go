package js

type Watcher struct{}

func NewWatcher() (*Watcher, error) {
	return nil, nil
}

func (self *Watcher) Add(path string) error {
	return nil
}

func (self *Watcher) Close() error {
	return nil
}

func (self *Environment) Watch(onChanged OnChangedFunc) error {
	return nil
}
