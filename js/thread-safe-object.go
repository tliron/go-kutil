package js

import (
	"sync"

	"github.com/dop251/goja"
)

type ThreadSafeObject struct {
	map_ map[string]goja.Value
	lock sync.Mutex
}

func NewThreadSafeObject() *ThreadSafeObject {
	return &ThreadSafeObject{
		map_: make(map[string]goja.Value),
	}
}

func (self *ThreadSafeObject) NewDynamicObject(runtime *goja.Runtime) *goja.Object {
	return runtime.NewDynamicObject(self)
}

// goja.DynamicObject interface
func (self *ThreadSafeObject) Get(key string) goja.Value {
	self.lock.Lock()
	defer self.lock.Unlock()

	value, _ := self.map_[key]
	return value
}

// goja.DynamicObject interface
func (self *ThreadSafeObject) Set(key string, value goja.Value) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.map_[key] = value
	return true
}

// goja.DynamicObject interface
func (self *ThreadSafeObject) Has(key string) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.map_[key]
	return ok
}

// goja.DynamicObject interface
func (self *ThreadSafeObject) Delete(key string) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.map_, key)
	return true
}

// goja.DynamicObject interface
func (self *ThreadSafeObject) Keys() []string {
	self.lock.Lock()
	defer self.lock.Unlock()

	keys := make([]string, 0, len(self.map_))
	for key := range self.map_ {
		keys = append(keys, key)
	}
	return keys
}
