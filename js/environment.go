package js

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/tliron/kutil/logging"
	urlpkg "github.com/tliron/kutil/url"
)

//
// Environment
//

type Environment struct {
	Runtime        *goja.Runtime
	URLContext     *urlpkg.Context
	Watcher        *Watcher
	Extensions     []Extension
	Modules        *goja.Object
	Precompile     PrecompileFunc
	CreateResolver CreateResolverFunc
	Log            logging.Logger

	exportsCache sync.Map
	programCache sync.Map
}

type PrecompileFunc func(url urlpkg.URL, script string, context *Context) (string, error)

type OnChangedFunc func(id string, module *Module)

func NewEnvironment(urlContext *urlpkg.Context) *Environment {
	self := Environment{
		Runtime:    goja.New(),
		URLContext: urlContext,
		CreateResolver: func(url urlpkg.URL, context *Context) ResolveFunc {
			var origins []urlpkg.URL
			if url != nil {
				origins = []urlpkg.URL{url.Origin()}
			}
			return func(id string) (urlpkg.URL, error) {
				/*if strings.IndexRune(id, '.') == -1 {
					id += ".js"
				}*/
				return urlpkg.NewValidURL(id, origins, urlContext)
			}
		},
		Log: log,
	}

	self.Modules = NewThreadSafeObject().NewDynamicObject(self.Runtime)

	self.Runtime.SetFieldNameMapper(CamelCaseMapper)

	return &self
}

func (self *Environment) Release() error {
	if self.Watcher != nil {
		if err := self.Watcher.Close(); err == nil {
			self.Watcher = nil
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

func (self *Environment) RequireID(id string) (*goja.Object, error) {
	return self.requireId(id, self.NewContext(nil, nil))
}

func (self *Environment) RequireURL(url urlpkg.URL) (*goja.Object, error) {
	return self.requireUrl(url, self.NewContext(url, nil))
}

func (self *Environment) requireId(id string, context *Context) (*goja.Object, error) {
	if url, err := context.Resolve(id); err == nil {
		self.AddModule(url, context.Module)
		return self.requireUrl(url, context)
	} else {
		return nil, err
	}
}

func (self *Environment) requireUrl(url urlpkg.URL, context *Context) (*goja.Object, error) {
	key := url.Key()

	// Try cache
	if exports, loaded := self.exportsCache.Load(key); loaded {
		// Cache hit
		return exports.(*goja.Object), nil
	} else {
		// Cache miss
		if exports, err := self.require_(url, context); err == nil {
			if exports_, loaded := self.exportsCache.LoadOrStore(key, exports); loaded {
				// Cache hit
				return exports_.(*goja.Object), nil
			} else {
				// Cache miss
				return exports, nil
			}
		} else {
			return nil, err
		}
	}
}

func (self *Environment) require_(url urlpkg.URL, context *Context) (*goja.Object, error) {
	context = self.NewContext(url, context)

	if program, err := self.compile(url, context); err == nil {
		if value, err := self.Runtime.RunProgram(program); err == nil {
			if call, ok := goja.AssertFunction(value); ok {
				// See: self.compile_ for arguments
				arguments := []goja.Value{
					context.Module.Exports,
					context.Module.Require,
					self.Runtime.ToValue(context.Module),
					self.Runtime.ToValue(context.Module.Filename),
					self.Runtime.ToValue(context.Module.Path),
				}

				arguments = append(arguments, context.Extensions...)

				if _, err := call(nil, arguments...); err == nil {
					return context.Module.Exports, nil
				} else {
					return nil, err
				}
			} else {
				// Should never happen
				return nil, fmt.Errorf("invalid module: %v", value)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Environment) compile(url urlpkg.URL, context *Context) (*goja.Program, error) {
	key := url.Key()

	// Try cache
	if program, loaded := self.programCache.Load(key); loaded {
		// Cache hit
		return program.(*goja.Program), nil
	} else {
		// Cache miss
		if program, err := self.compile_(url, context); err == nil {
			if program_, loaded := self.programCache.LoadOrStore(key, program); loaded {
				// Cache hit
				return program_.(*goja.Program), nil
			} else {
				// Cache miss
				return program, nil
			}
		} else {
			return nil, err
		}
	}
}

func (self *Environment) compile_(url urlpkg.URL, context *Context) (*goja.Program, error) {
	if script, err := urlpkg.ReadString(url); err == nil {
		// Precompile
		if self.Precompile != nil {
			if script, err = self.Precompile(url, script, context); err != nil {
				return nil, err
			}
		}

		// See: https://nodejs.org/api/modules.html#modules_the_module_wrapper
		var builder strings.Builder
		builder.WriteString("(function(exports, require, module, __filename, __dirname")
		for _, extension := range self.Extensions {
			builder.WriteString(", ")
			builder.WriteString(extension.Name)
		}
		builder.WriteString(") {")
		builder.WriteString(script)
		builder.WriteString("\n});")
		script = builder.String()
		//log.Infof("%s", script)

		return goja.Compile(url.String(), script, true)
	} else {
		return nil, err
	}
}
