package url

import (
	"fmt"
	"io"
	"net/http"
	neturlpkg "net/url"
	"path"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// NetworkURL
//

type NetworkURL struct {
	URL *neturlpkg.URL

	string_ string
	context *Context
}

func NewNetworkURL(neturl *neturlpkg.URL, context *Context) *NetworkURL {
	if context == nil {
		context = NewContext()
	}

	return &NetworkURL{
		URL:     neturl,
		string_: neturl.String(),
		context: context,
	}
}

func NewValidNetworkURL(neturl *neturlpkg.URL, context *Context) (*NetworkURL, error) {
	string_ := neturl.String()
	if response, err := http.Get(string_); err == nil {
		response.Body.Close()
		if response.StatusCode == http.StatusOK {
			if context == nil {
				context = NewContext()
			}

			return &NetworkURL{
				URL:     neturl,
				string_: string_,
				context: context,
			}, nil
		} else {
			return nil, fmt.Errorf("HTTP status: %s", response.Status)
		}
	} else {
		return nil, err
	}
}

func NewValidRelativeNetworkURL(path string, origin *NetworkURL) (*NetworkURL, error) {
	if neturl, err := neturlpkg.Parse(path); err == nil {
		neturl = origin.URL.ResolveReference(neturl)
		return NewValidNetworkURL(neturl, origin.context)
	} else {
		return nil, err
	}
}

// URL interface
// fmt.Stringer interface
func (self *NetworkURL) String() string {
	return self.Key()
}

// URL interface
func (self *NetworkURL) Format() string {
	format := self.URL.Query().Get("format")
	if format != "" {
		return format
	} else {
		return GetFormat(self.URL.Path)
	}
}

// URL interface
func (self *NetworkURL) Origin() URL {
	url := *self
	url.URL.Path = path.Dir(url.URL.Path)
	return &url
}

// URL interface
func (self *NetworkURL) Relative(path string) URL {
	if neturl, err := neturlpkg.Parse(path); err == nil {
		return NewNetworkURL(self.URL.ResolveReference(neturl), self.context)
	} else {
		return nil
	}
}

// URL interface
func (self *NetworkURL) Key() string {
	return self.string_
}

// URL interface
func (self *NetworkURL) Open() (io.ReadCloser, error) {
	if response, err := http.Get(self.string_); err == nil {
		if response.StatusCode == http.StatusOK {
			return response.Body, nil
		} else {
			response.Body.Close()
			return nil, fmt.Errorf("HTTP status: %s", response.Status)
		}
	} else {
		return nil, err
	}
}

// URL interface
func (self *NetworkURL) Context() *Context {
	return self.context
}
