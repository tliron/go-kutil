package util

import (
	"sync"
)

//
// Promise
//

type Promise chan struct{}

func NewPromise() Promise {
	return make(Promise)
}

func (self Promise) Release() {
	close(self)
}

func (self Promise) Wait() {
	<-self
}

//
// CoordinatedWork
//

type CoordinatedWork struct {
	sync.Map
}

func NewCoordinatedWork() *CoordinatedWork {
	return &CoordinatedWork{}
}

func (self *CoordinatedWork) Start(key string) (Promise, bool) {
	promise := NewPromise()
	if existing, loaded := self.LoadOrStore(key, promise); !loaded {
		return promise, true
	} else {
		promise = existing.(Promise)
		promise.Wait()
		return nil, false
	}
}
