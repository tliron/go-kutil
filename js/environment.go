package js

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/tliron/kutil/fswatch"
	"github.com/tliron/kutil/logging"
	urlpkg "github.com/tliron/kutil/url"
)

//
// Environment
//

type Environment struct {
	Runtime        *goja.Runtime
	URLContext     *urlpkg.Context
	Extensions     []Extension
	Modules        *goja.Object
	Precompile     PrecompileFunc
	CreateResolver CreateResolverFunc
	Log            logging.Logger
	Lock           sync.Mutex

	watcher      *fswatch.Watcher
	watcherLock  sync.Mutex
	exportsCache sync.Map
	programCache *sync.Map
}

type PrecompileFunc func(url urlpkg.URL, script string, context *Context) (string, error)

type OnChangedFunc func(id string, module *Module)

func NewEnvironment(urlContext *urlpkg.Context) *Environment {
	self := Environment{
		Runtime:        goja.New(),
		URLContext:     urlContext,
		CreateResolver: NewDefaultResolverCreator(urlContext, "js"),
		Log:            log,
		programCache:   new(sync.Map),
	}

	self.Modules = NewThreadSafeObject().NewDynamicObject(self.Runtime)

	self.Runtime.SetFieldNameMapper(CamelCaseMapper)

	return &self
}

func (self *Environment) NewChild() *Environment {
	environment := NewEnvironment(self.URLContext)
	environment.watcher = self.watcher
	environment.Extensions = self.Extensions
	environment.Precompile = self.Precompile
	environment.CreateResolver = self.CreateResolver
	environment.Log = self.Log
	environment.programCache = self.programCache
	return environment
}

func (self *Environment) StartWatcher(onChanged OnChangedFunc) error {
	self.watcherLock.Lock()
	defer self.watcherLock.Unlock()

	if watcher, err := fswatch.NewWatcher(self.URLContext); err == nil {
		self.watcher = watcher
		watcher.Start(func(fileUrl *urlpkg.FileURL) {
			self.Lock.Lock()
			id := fileUrl.Key()
			var module *Module
			if module_ := self.Modules.Get(id); module_ != nil {
				module = module_.Export().(*Module)
			}
			self.Lock.Unlock()
			onChanged(id, module)
		})
		return nil
	} else {
		return err
	}
}

func (self *Environment) Watch(path string) error {
	self.watcherLock.Lock()
	defer self.watcherLock.Unlock()

	if self.watcher != nil {
		return self.watcher.Add(path)
	} else {
		return nil
	}
}

func (self *Environment) Release() error {
	self.watcherLock.Lock()
	defer self.watcherLock.Unlock()

	if self.watcher != nil {
		if err := self.watcher.Close(); err == nil {
			self.watcher = nil
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

func (self *Environment) Call(function JavaScriptFunc, arguments ...interface{}) interface{} {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	return Call(self.Runtime, function, arguments...)
}

func (self *Environment) ClearCache() {
	self.exportsCache.Range(func(key interface{}, value interface{}) bool {
		self.exportsCache.Delete(key)
		return true
	})
	self.programCache.Range(func(key interface{}, value interface{}) bool {
		self.programCache.Delete(key)
		return true
	})
	self.Modules = NewThreadSafeObject().NewDynamicObject(self.Runtime)
}

func (self *Environment) RequireID(id string) (*goja.Object, error) {
	return self.requireId(id, self.NewContext(nil, nil))
}

func (self *Environment) RequireURL(url urlpkg.URL) (*goja.Object, error) {
	return self.cachedRequire(url, self.NewContext(url, nil))
}

func (self *Environment) requireId(id string, context *Context) (*goja.Object, error) {
	if url, err := context.Resolve(id, false); err == nil {
		self.AddModule(url, context.Module)
		return self.cachedRequire(url, context)
	} else {
		return nil, err
	}
}

func (self *Environment) cachedRequire(url urlpkg.URL, context *Context) (*goja.Object, error) {
	key := url.Key()

	// Try cache
	if exports, loaded := self.exportsCache.Load(key); loaded {
		// Cache hit
		return exports.(*goja.Object), nil
	} else {
		// Cache miss
		if exports, err := self.require(url, context); err == nil {
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

func (self *Environment) require(url urlpkg.URL, context *Context) (*goja.Object, error) {
	// Create a child context
	context = self.NewContext(url, context)

	if program, err := self.cachedCompile(url, context); err == nil {
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

func (self *Environment) cachedCompile(url urlpkg.URL, context *Context) (*goja.Program, error) {
	key := url.Key()

	// Try cache
	if program, loaded := self.programCache.Load(key); loaded {
		// Cache hit
		return program.(*goja.Program), nil
	} else {
		// Cache miss
		if program, err := self.compile(url, context); err == nil {
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

func (self *Environment) compile(url urlpkg.URL, context *Context) (*goja.Program, error) {
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
