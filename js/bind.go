package js

import (
	"github.com/dop251/goja"
	urlpkg "github.com/tliron/kutil/url"
)

//
// Bind
//

type Bind interface {
	Unbind() (interface{}, *Context, error)
}

//
// EarlyBind
//

type EarlyBind struct {
	Value   interface{}
	Context *Context
}

// Bind interface
func (self EarlyBind) Unbind() (interface{}, *Context, error) {
	return self.Value, self.Context, nil
}

// CreateExtensionFunc signature
func CreateEarlyBindExtension(context *Context) goja.Value {
	return context.Environment.Runtime.ToValue(func(id string, exportName string) (goja.Value, error) {
		if url, err := context.Resolve(id, false); err == nil {
			childEnvironment := context.Environment.NewChild()
			childContext := childEnvironment.NewContext(url, nil)
			if exports, err := childEnvironment.cachedRequire(url, childContext); err == nil {
				var value interface{}

				if exportName == "" {
					value = exports.Export()
				} else {
					value = exports.Get(exportName).Export()
				}

				return context.Environment.Runtime.ToValue(EarlyBind{
					Value:   value,
					Context: childContext,
				}), nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	})
}

//
// LateBind
//

type LateBind struct {
	URL        urlpkg.URL
	ExportName string
	Context    *Context
}

// Bind interface
func (self LateBind) Unbind() (interface{}, *Context, error) {
	childEnvironment := self.Context.Environment.NewChild()
	childContext := childEnvironment.NewContext(self.URL, nil)
	if exports, err := childEnvironment.cachedRequire(self.URL, childContext); err == nil {
		var value interface{}

		if self.ExportName == "" {
			value = exports.Export()
		} else {
			value = exports.Get(self.ExportName).Export()
		}

		return value, childContext, nil
	} else {
		return nil, nil, err
	}
}

// CreateExtensionFunc signature
func CreateLateBindExtension(context *Context) goja.Value {
	return context.Environment.Runtime.ToValue(func(id string, exportName string) (goja.Value, error) {
		if url, err := context.Resolve(id, false); err == nil {
			return context.Environment.Runtime.ToValue(LateBind{
				URL:        url,
				ExportName: exportName,
				Context:    context,
			}), nil
		} else {
			return nil, err
		}
	})
}
