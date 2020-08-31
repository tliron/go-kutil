package util

import (
	"sync"
)

var locks sync.Map

func GetLock(pointer interface{}) *sync.RWMutex {
	lock := new(sync.RWMutex)
	if existing, loaded := locks.LoadOrStore(pointer, lock); loaded {
		return existing.(*sync.RWMutex)
	} else {
		return lock
	}
}

func ResetLocks() {
	// See: https://stackoverflow.com/a/49355523
	locks.Range(func(key interface{}, value interface{}) bool {
		locks.Delete(key)
		return true
	})
}
