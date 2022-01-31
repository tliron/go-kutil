package util

import (
	"fmt"
	"sync"

	"github.com/sasha-s/go-deadlock"
)

func init() {
	//deadlock.Opts.DisableLockOrderDetection = true
}

//
// RWLocker
//

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
	RLocker() sync.Locker
}

type LockType int

const (
	DEFAULT_LOCK = LockType(0)
	DEBUG_LOCK   = LockType(1)
	MOCK_LOCK    = LockType(2)
)

func NewRWLocker(type_ LockType) RWLocker {
	switch type_ {
	case DEFAULT_LOCK:
		return NewDefaultRWLocker()
	case DEBUG_LOCK:
		return NewDebugRWLocker()
	case MOCK_LOCK:
		return NewMockRWLocker()
	default:
		panic(fmt.Sprintf("unsupported lock type: %d", type_))
	}
}

//
// DefaultRWLocker
//

func NewDefaultRWLocker() RWLocker {
	return new(sync.RWMutex)
}

//
// DebugRWLocker
//

func NewDebugRWLocker() RWLocker {
	return new(deadlock.RWMutex)
}

//
// MockLocker
//

type MockLocker struct{}

func NewMockLocker() sync.Locker {
	return MockLocker{}
}

// sync.Locker interface
func (self MockLocker) Lock() {}

// sync.Locker interface
func (self MockLocker) Unlock() {}

//
// MockRWLocker
//

type MockRWLocker struct {
	MockLocker
}

func NewMockRWLocker() RWLocker {
	return MockRWLocker{}
}

// RWLocker interface
func (self MockRWLocker) RLock() {}

// RWLocker interface
func (self MockRWLocker) RUnlock() {}

// RWLocker interface
func (self MockRWLocker) RLocker() sync.Locker {
	return self
}

//
// LockableEntity
//

type LockableEntity interface {
	GetEntityLock() RWLocker
}

func GetEntityLock(entity interface{}) RWLocker {
	if lockable, ok := entity.(LockableEntity); ok {
		//fmt.Printf("entity: %T %v\n", entity, entity)
		return lockable.GetEntityLock()
	} else {
		return nil
	}
}

//
// Ad-hoc locks
//

var adHocLocks sync.Map

// Warning: Addresses can be re-used after the resource if freed.
// This facililty should only be used for long-lived objects.
func GetAdHocLock(pointer interface{}, type_ LockType) RWLocker {
	if pointer == nil {
		panic("no ad-hoc lock for nil")
	}

	if lock, ok := adHocLocks.Load(pointer); ok {
		return lock.(RWLocker)
	} else {
		lock := NewRWLocker(type_)
		if existing, loaded := adHocLocks.LoadOrStore(pointer, lock); loaded {
			return existing.(RWLocker)
		} else {
			return lock
		}
	}
}

func ResetAdHocLocks() {
	// See: https://stackoverflow.com/a/49355523
	adHocLocks.Range(func(key interface{}, value interface{}) bool {
		adHocLocks.Delete(key)
		return true
	})
}
