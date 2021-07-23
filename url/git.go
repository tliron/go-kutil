package url

import (
	"fmt"
	"io"
	neturlpkg "net/url"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/tliron/kutil/util"
)

//
// GitURL
//

type GitURL struct {
	Path          string
	RepositoryURL string
	Reference     string
	Username      string
	Password      string

	clonePath string
	context   *Context
}

func NewGitURL(path string, repositoryUrl string, context *Context) *GitURL {
	// Must be absolute
	path = strings.TrimLeft(path, "/")

	var self = GitURL{
		Path:    path,
		context: context,
	}

	if neturl, err := neturlpkg.Parse(repositoryUrl); err == nil {
		if neturl.User != nil {
			self.Username = neturl.User.Username()
			if password, ok := neturl.User.Password(); ok {
				self.Password = password
			}
			// Don't store user info
			neturl.User = nil
		}
		self.Reference = neturl.Fragment
		self.RepositoryURL = neturl.String()
	} else {
		self.RepositoryURL = repositoryUrl
	}

	return &self
}

func NewValidGitURL(path string, repositoryUrl string, context *Context) (*GitURL, error) {
	self := NewGitURL(path, repositoryUrl, context)
	if _, err := self.OpenRepository(); err == nil {
		path := filepath.Join(self.clonePath, self.Path)
		if _, err := os.Stat(path); err == nil {
			return self, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewValidRelativeGitURL(path string, origin *GitURL) (*GitURL, error) {
	self := origin.Relative(path).(*GitURL)
	if _, err := self.OpenRepository(); err == nil {
		path_ := filepath.Join(self.clonePath, self.Path)
		if _, err := os.Stat(path_); err == nil {
			return self, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ParseGitURL(url string, context *Context) (*GitURL, error) {
	if repositoryUrl, path, err := parseGitURL(url); err == nil {
		return NewGitURL(path, repositoryUrl, context), nil
	} else {
		return nil, err
	}
}

func ParseValidGitURL(url string, context *Context) (*GitURL, error) {
	if repositoryUrl, path, err := parseGitURL(url); err == nil {
		return NewValidGitURL(path, repositoryUrl, context)
	} else {
		return nil, err
	}
}

// URL interface
// fmt.Stringer interface
func (self *GitURL) String() string {
	return self.Key()
}

// URL interface
func (self *GitURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *GitURL) Origin() URL {
	path := pathpkg.Dir(self.Path)
	if path != "/" {
		path += "/"
	}

	return &GitURL{
		Path:          path,
		RepositoryURL: self.RepositoryURL,
		clonePath:     self.clonePath,
		context:       self.context,
	}
}

// URL interface
func (self *GitURL) Relative(path string) URL {
	return &GitURL{
		Path:          pathpkg.Join(self.Path, path),
		RepositoryURL: self.RepositoryURL,
		clonePath:     self.clonePath,
		context:       self.context,
	}
}

// URL interface
func (self *GitURL) Key() string {
	return fmt.Sprintf("git:%s!/%s", self.RepositoryURL, self.Path)
}

// URL interface
func (self *GitURL) Open() (io.ReadCloser, error) {
	if _, err := self.OpenRepository(); err == nil {
		path := filepath.Join(self.clonePath, self.Path)
		return os.Open(path)
	} else {
		return nil, err
	}
}

// URL interface
func (self *GitURL) Context() *Context {
	return self.context
}

func (self *GitURL) OpenRepository() (*git.Repository, error) {
	if self.clonePath != "" {
		return self.openRepository(false)
	} else {
		key := self.Key()

		// Note: this will lock for the entire clone duration!
		self.context.lock.Lock()
		defer self.context.lock.Unlock()

		if self.context.dirs != nil {
			// Already cloned?
			if clonePath, ok := self.context.dirs[key]; ok {
				self.clonePath = clonePath
				return self.openRepository(false)
			}
		}

		temporaryPathPattern := fmt.Sprintf("kutil-%s-*", util.SanitizeFilename(key))
		if clonePath, err := os.MkdirTemp("", temporaryPathPattern); err == nil {
			if self.context.dirs == nil {
				self.context.dirs = make(map[string]string)
			}
			self.context.dirs[key] = clonePath

			// Clone
			if repository, err := git.PlainClone(clonePath, false, &git.CloneOptions{
				URL:  self.RepositoryURL,
				Auth: self.getAuth(),
			}); err == nil {
				if reference, err := self.findReference(repository); err == nil {
					if reference != nil {
						// Checkout
						if workTree, err := repository.Worktree(); err == nil {
							if err := workTree.Checkout(&git.CheckoutOptions{
								Branch: reference.Name(),
							}); err != nil {
								os.RemoveAll(clonePath)
								return nil, err
							}
						} else {
							os.RemoveAll(clonePath)
							return nil, err
						}
					}
				} else {
					os.RemoveAll(clonePath)
					return nil, err
				}

				self.clonePath = clonePath
				return repository, nil
			} else {
				os.RemoveAll(clonePath)
				return nil, err
			}
		} else {
			return nil, err
		}
	}
}

func (self *GitURL) openRepository(pull bool) (*git.Repository, error) {
	if repository, err := git.PlainOpen(self.clonePath); err == nil {
		if pull {
			if err := self.pullRepository(repository); err != nil {
				return nil, err
			}
		}

		return repository, nil
	} else {
		return nil, err
	}
}

func (self *GitURL) pullRepository(repository *git.Repository) error {
	if workTree, err := repository.Worktree(); err == nil {
		if err := workTree.Pull(&git.PullOptions{
			Auth: self.getAuth(),
		}); (err == nil) || (err == git.NoErrAlreadyUpToDate) {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *GitURL) findReference(repository *git.Repository) (*plumbing.Reference, error) {
	if self.Reference != "" {
		if iter, err := repository.References(); err == nil {
			defer iter.Close()
			for {
				if reference, err := iter.Next(); err == nil {
					name := reference.Name()
					if name.Short() == self.Reference {
						return reference, nil
					} else if name.String() == self.Reference {
						return reference, nil
					}
				} else if err == io.EOF {
					return nil, fmt.Errorf("reference %q not found in git repository: %s", self.Reference, self.RepositoryURL)
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		return nil, nil
	}
}

func (self *GitURL) getAuth() transport.AuthMethod {
	// TODO: what about non-HTTP transports, like ssh?
	if self.Username != "" {
		return &http.BasicAuth{
			Username: self.Username,
			Password: self.Password,
		}
	} else {
		return nil
	}
}

func parseGitURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "git:") {
		if split := strings.Split(url[4:], "!"); len(split) == 2 {
			return split[0], split[1], nil
		} else {
			return "", "", fmt.Errorf("malformed \"git:\" URL: %s", url)
		}
	} else {
		return "", "", fmt.Errorf("not a \"git:\" URL: %s", url)
	}
}
