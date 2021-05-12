package js

import (
	"github.com/dop251/goja"
)

//
// Hook
//
// Struct rather than function so that goja won't wrap it
//

type Hook struct {
	Callable goja.Callable
	Call     func(this interface{}, args ...interface{}) (interface{}, error)
}

func NewHook(callable goja.Callable, runtime *goja.Runtime) *Hook {
	return &Hook{
		Callable: callable,
		Call: func(this interface{}, args ...interface{}) (interface{}, error) {
			args_ := make([]goja.Value, len(args))
			for index, arg := range args {
				args_[index] = runtime.ToValue(arg)
			}

			if r, err := callable(runtime.ToValue(this), args_...); err == nil {
				return r.Export(), nil
			} else {
				return nil, err
			}
		},
	}
}
