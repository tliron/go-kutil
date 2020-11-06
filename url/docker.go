package url

import (
	"fmt"
	"io"
	neturlpkg "net/url"
	"path"

	"github.com/google/go-containerregistry/pkg/authn"
	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

//
// DockerURL
//

type DockerURL struct {
	URL *neturlpkg.URL

	string_ string
	context *Context
}

func NewDockerURL(neturl *neturlpkg.URL, context *Context) *DockerURL {
	if context == nil {
		context = NewContext()
	}

	return &DockerURL{
		URL:     neturl,
		string_: neturl.String(),
		context: context,
	}
}

func NewValidDockerURL(neturl *neturlpkg.URL, context *Context) (*DockerURL, error) {
	if (neturl.Scheme != "docker") && (neturl.Scheme != "") {
		return nil, fmt.Errorf("not a docker URL: %s", neturl.String())
	}

	// TODO

	return NewDockerURL(neturl, context), nil
}

// URL interface
// fmt.Stringer interface
func (self *DockerURL) String() string {
	return self.Key()
}

// URL interface
func (self *DockerURL) Format() string {
	format := self.URL.Query().Get("format")
	if format != "" {
		return format
	} else {
		return GetFormat(self.URL.Path)
	}
}

// URL interface
func (self *DockerURL) Origin() URL {
	url := *self
	url.URL.Path = path.Dir(url.URL.Path)
	return &url
}

// URL interface
func (self *DockerURL) Relative(path string) URL {
	if neturl, err := neturlpkg.Parse(path); err == nil {
		return NewDockerURL(self.URL.ResolveReference(neturl), self.context)
	} else {
		return nil
	}
}

// URL interface
func (self *DockerURL) Key() string {
	return self.string_
}

// URL interface
func (self *DockerURL) Open() (io.ReadCloser, error) {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		if err := self.WriteLayer(pipeWriter); err == nil {
			pipeWriter.Close()
		} else {
			pipeWriter.CloseWithError(err)
		}
	}()

	return pipeReader, nil
}

// URL interface
func (self *DockerURL) Context() *Context {
	return self.context
}

func (self *DockerURL) WriteTarball(writer io.Writer) error {
	url := fmt.Sprintf("%s%s", self.URL.Host, self.URL.Path)
	if tag, err := namepkg.NewTag(url); err == nil {
		if image, err := remote.Image(tag, self.RemoteOptions()...); err == nil {
			return tarball.Write(tag, image, writer)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *DockerURL) WriteLayer(writer io.Writer) error {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		if err := self.WriteTarball(pipeWriter); err == nil {
			pipeWriter.Close()
		} else {
			pipeWriter.CloseWithError(err)
		}
	}()

	decoder := NewContainerImageLayerDecoder(pipeReader)
	if _, err := io.Copy(writer, decoder.Decode()); err == nil {
		return nil
	} else {
		return err
	}
}

func (self *DockerURL) RemoteOptions() []remote.Option {
	var options []remote.Option

	if httpRoundTripper := self.context.GetHTTPRoundTripper(self.URL.Host); httpRoundTripper != nil {
		options = append(options, remote.WithTransport(httpRoundTripper))
	}

	if credentials := self.context.GetCredentials(self.URL.Host); credentials != nil {
		authenticator := authn.FromConfig(authn.AuthConfig{
			Username:      credentials.Username,
			Password:      credentials.Password,
			RegistryToken: credentials.Token,
		})
		options = append(options, remote.WithAuth(authenticator))
	}

	return options
}
