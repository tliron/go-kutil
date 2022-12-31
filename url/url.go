package url

import (
	"errors"
	"fmt"
	"io"
	neturlpkg "net/url"
	pathpkg "path"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// URL
//

type URL interface {
	String() string
	Format() string // yaml|json|xml
	Origin() URL    // base dir, is not necessarily a valid URL
	Relative(path string) URL
	Key() string // for maps
	Open() (io.ReadCloser, error)
	Context() *Context
}

func NewURL(url string, context *Context) (URL, error) {
	if context != nil {
		if url_, ok := context.GetMapping(url); ok {
			url = url_
		}
	}

	neturl, err := neturlpkg.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("unsupported URL format: %s", url)
	} else {
		switch neturl.Scheme {
		case "file":
			return NewFileURL(neturl.Path, context), nil

		case "http", "https":
			// Go's "net/http" only handles "http:" and "https:"
			return NewNetworkURL(neturl, context), nil

		case "tar":
			return ParseTarballURL(url, context)

		case "zip":
			return ParseZipURL(url, context)

		case "git":
			return ParseGitURL(url, context)

		case "docker":
			return NewDockerURL(neturl, context), nil

		case "internal":
			return NewInternalURL(url[9:], context), nil

		case "":
			return NewFileURL(url, context), nil
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func NewValidURL(url string, origins []URL, context *Context) (URL, error) {
	if context != nil {
		if url_, ok := context.GetMapping(url); ok {
			url = url_
		}
	}

	neturl, err := neturlpkg.ParseRequestURI(url)
	if err != nil {
		// Malformed URL, so it might be a relative path
		return newValidRelativeURL(url, origins, context, false)
	} else {
		switch neturl.Scheme {
		case "file":
			// They're rarely used, but relative "file:" URLs are possible
			return newValidRelativeURL(neturl.Path, origins, context, true)

		case "http", "https":
			// Go's "net/http" package only handles "http:" and "https:"
			return NewValidNetworkURL(neturl, context)

		case "tar":
			return ParseValidTarballURL(url, context)

		case "zip":
			return ParseValidZipURL(url, context)

		case "git":
			return ParseValidGitURL(url, context)

		case "docker":
			return NewValidDockerURL(neturl, context)

		case "internal":
			return NewValidInternalURL(url[9:], context)

		case "":
			return newValidRelativeURL(url, origins, context, false)
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func newValidRelativeURL(path string, origins []URL, context *Context, onlyFileURLs bool) (URL, error) {
	// Absolute file path?
	if pathpkg.IsAbs(path) {
		url, err := NewValidFileURL(path, context)
		if err != nil {
			return nil, err
		}
		return url, nil
	} else {
		// Try relative to origins
		for _, origin := range origins {
			var url URL
			err := errors.New("")

			switch origin_ := origin.(type) {
			case *FileURL:
				url, err = NewValidRelativeFileURL(path, origin_)

			case *NetworkURL:
				if !onlyFileURLs {
					url, err = NewValidRelativeNetworkURL(path, origin_)
				}

			case *TarballURL:
				if !onlyFileURLs {
					url, err = NewValidRelativeTarballURL(path, origin_)
				}

			case *ZipURL:
				if !onlyFileURLs {
					url, err = NewValidRelativeZipURL(path, origin_)
				}

			case *GitURL:
				if !onlyFileURLs {
					url, err = NewValidRelativeGitURL(path, origin_)
				}

			case *InternalURL:
				if !onlyFileURLs {
					url, err = NewValidRelativeInternalURL(path, origin_)
				}
			}

			if err == nil {
				return url, nil
			}
		}

		// Try file relative to current directory
		url, err := NewValidFileURL(path, context)
		if err != nil {
			return nil, NewNotFoundf("URL not found: %s", path)
		}

		return url, nil
	}
}
