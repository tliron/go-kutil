package js

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	urlpkg "github.com/tliron/kutil/url"
)

type ResolveFunc func(id string) (urlpkg.URL, error)

type CreateResolverFunc func(url urlpkg.URL) ResolveFunc

type PrecompileFunc func(url urlpkg.URL, script string, resolve ResolveFunc) (string, error)

type CreateExtensionFunc func(environment *Environment, resolve ResolveFunc) goja.Value

type Extension struct {
	Name   string
	Create CreateExtensionFunc
}

//
// Environment
//

type Environment struct {
	Runtime        *goja.Runtime
	URLContext     *urlpkg.Context
	Extensions     []Extension
	Precompile     PrecompileFunc
	CreateResolver CreateResolverFunc

	exportsCache sync.Map
	programCache sync.Map
}

func NewEnvironment(urlContext *urlpkg.Context) *Environment {
	self := Environment{
		Runtime:    goja.New(),
		URLContext: urlContext,
		CreateResolver: func(url urlpkg.URL) ResolveFunc {
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
	}

	self.Runtime.SetFieldNameMapper(CamelCaseMapper)

	return &self
}

func (self *Environment) RequireURL(url urlpkg.URL) (*goja.Object, error) {
	return self.require(url, self.CreateResolver(url))
}

func (self *Environment) Require(id string) (*goja.Object, error) {
	return self.ResolveAndRequire(id, self.CreateResolver(nil))
}

func (self *Environment) ResolveAndRequire(id string, resolve ResolveFunc) (*goja.Object, error) {
	if url, err := resolve(id); err == nil {
		return self.require(url, resolve)
	} else {
		return nil, err
	}
}

func (self *Environment) require(url urlpkg.URL, resolve ResolveFunc) (*goja.Object, error) {
	key := url.Key()

	// Try cache
	if exports, loaded := self.exportsCache.Load(key); loaded {
		// Cache hit
		return exports.(*goja.Object), nil
	} else {
		// Cache miss
		if exports, err := self.require_(url, resolve); err == nil {
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

func (self *Environment) require_(url urlpkg.URL, resolve ResolveFunc) (*goja.Object, error) {
	if program, err := self.compile(url, resolve); err == nil {
		if value, err := self.Runtime.RunProgram(program); err == nil {
			if call, ok := goja.AssertFunction(value); ok {
				module, extensions := self.newModule(url)

				// See: self.compile_ for arguments
				arguments := []goja.Value{module.Get("exports"), module.Get("require"), module, module.Get("filename"), module.Get("path")}
				arguments = append(arguments, extensions...)
				if _, err := call(nil, arguments...); err == nil {
					return module.Get("exports").(*goja.Object), nil
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

func (self *Environment) compile(url urlpkg.URL, resolve ResolveFunc) (*goja.Program, error) {
	key := url.Key()

	// Try cache
	if program, loaded := self.programCache.Load(key); loaded {
		// Cache hit
		return program.(*goja.Program), nil
	} else {
		// Cache miss
		if program, err := self.compile_(url, resolve); err == nil {
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

func (self *Environment) compile_(url urlpkg.URL, resolve ResolveFunc) (*goja.Program, error) {
	if script, err := urlpkg.ReadString(url); err == nil {
		// Precompile
		if self.Precompile != nil {
			if script, err = self.Precompile(url, script, resolve); err != nil {
				return nil, err
			}
		}

		// See: https://nodejs.org/api/modules.html#modules_the_module_wrapper
		var extensions string
		for _, extension := range self.Extensions {
			extensions += ", " + extension.Name
		}
		script = "(function(exports, require, module, __filename, __dirname" + extensions + ") {" + script + "\n});"
		//log.Infof("%s", script)

		return goja.Compile(url.String(), script, true)
	} else {
		return nil, err
	}
}

func (self *Environment) newModule(url urlpkg.URL) (*goja.Object, []goja.Value) {
	resolve := self.CreateResolver(url)

	var filename string
	var dirname string
	if url_, ok := url.(*urlpkg.FileURL); ok {
		filename = url_.Path
		if origin, ok := url_.Origin().(*urlpkg.FileURL); ok {
			dirname = origin.Path
		}
	}

	// See: https://nodejs.org/api/modules.html#modules_the_module_object
	module := self.Runtime.NewObject()
	module.Set("id", url.Key())
	module.Set("exports", self.Runtime.NewObject())
	module.Set("filename", filename)
	module.Set("path", dirname)

	// See: https://nodejs.org/api/modules.html#modules_require_id
	require := self.Runtime.ToValue(func(id string) (goja.Value, error) {
		return self.ResolveAndRequire(id, resolve)
	}).(*goja.Object)
	require.Set("cache", nil)
	require.Set("main", module)
	require.Set("resolve", func(id string, options *goja.Object) (string, error) {
		if url, err := resolve(id); err == nil {
			return url.String(), nil
		} else {
			return "", err
		}
	})
	module.Set("require", require)

	var extensions []goja.Value
	for _, extension := range self.Extensions {
		extensions = append(extensions, extension.Create(self, resolve))
	}

	return module, extensions
}
