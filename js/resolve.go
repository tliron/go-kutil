package js

import (
	"path/filepath"

	"github.com/tliron/exturl"
)

type ResolveFunc func(id string, raw bool) (exturl.URL, error)

type CreateResolverFunc func(url exturl.URL, context *Context) ResolveFunc

func NewDefaultResolverCreator(urlContext *exturl.Context, path []exturl.URL, defaultExtension string) CreateResolverFunc {
	return func(url exturl.URL, context *Context) ResolveFunc {
		var origins []exturl.URL

		if url != nil {
			origins = append([]exturl.URL{url.Origin()}, path...)
		} else {
			origins = path
		}

		if defaultExtension == "" {
			return func(id string, raw bool) (exturl.URL, error) {
				return exturl.NewValidURL(id, origins, urlContext)
			}
		} else {
			defaultExtension_ := "." + defaultExtension
			return func(id string, raw bool) (exturl.URL, error) {
				if !raw {
					if filepath.Ext(id) == "" {
						id += defaultExtension_
					}
				}

				return exturl.NewValidURL(id, origins, urlContext)
			}
		}
	}
}
