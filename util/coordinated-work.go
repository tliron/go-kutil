package util

import (
	contextpkg "context"
	"sync"
)

//
// CoordinatedWork
//

type CoordinatedWork struct {
	sync.Map
}

func NewCoordinatedWork() *CoordinatedWork {
	return new(CoordinatedWork)
}

func (self *CoordinatedWork) Start(context contextpkg.Context, key string) (Promise, bool) {
	promise := NewPromise()
	if existing, loaded := self.LoadOrStore(key, promise); !loaded {
		return promise, true
	} else {
		promise = existing.(Promise)
		promise.Wait(context)
		return nil, false
	}
}
