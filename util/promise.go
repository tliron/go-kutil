package util

import contextpkg "context"

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

func (self Promise) Wait(context contextpkg.Context) error {
	select {
	case <-context.Done():
		return context.Err()
	case <-self:
		return nil
	}
}
