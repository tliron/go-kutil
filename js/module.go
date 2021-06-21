package js

import (
	"github.com/dop251/goja"
	urlpkg "github.com/tliron/kutil/url"
)

//
// Module
//

type Module struct {
	Id           string
	Children     []*Module
	Filename     string
	Path         string
	Paths        []string
	Exports      *goja.Object
	Require      *goja.Object
	IsPreloading bool
	Loaded       bool
}

func (self *Environment) NewModule() *Module {
	var path []string
	for _, url := range self.Path {
		path = append(path, url.String())
	}

	return &Module{
		Paths:        path,
		Exports:      self.Runtime.NewObject(),
		IsPreloading: true,
	}
}

func (self *Environment) AddModule(url urlpkg.URL, module *Module) {
	module.Id = url.Key()
	module.IsPreloading = false
	module.Loaded = true
	if fileUrl, ok := url.(*urlpkg.FileURL); ok {
		module.Filename = fileUrl.Path
		if fileOrigin, ok := fileUrl.Origin().(*urlpkg.FileURL); ok {
			module.Path = fileOrigin.Path
		}

		if err := self.Watch(module.Filename); err != nil {
			self.Log.Errorf("%s", err.Error())
		}
	}

	self.Modules.Set(module.Id, module)
}
