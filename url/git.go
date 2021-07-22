package url

import (
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/tliron/kutil/util"
)

//
// GitURL
//

type GitURL struct {
	Path          string
	RepositoryURL string

	clonePath string
	context   *Context
}

func NewGitURL(path string, repositoryUrl string, context *Context) *GitURL {
	// Must be absolute
	path = strings.TrimLeft(path, "/")

	return &GitURL{
		Path:          path,
		RepositoryURL: repositoryUrl,
		context:       context,
	}
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
		return openAndPullGitRepository(self.clonePath, self.RepositoryURL)
	} else {
		key := self.Key()

		self.context.lock.Lock()
		defer self.context.lock.Unlock()

		if self.context.dirs != nil {
			if clonePath, ok := self.context.dirs[key]; ok {
				self.clonePath = clonePath
				return openAndPullGitRepository(self.clonePath, self.RepositoryURL)
			}
		}

		temporaryPathPattern := fmt.Sprintf("kutil-%s-*", util.SanitizeFilename(key))
		if clonePath, err := os.MkdirTemp("", temporaryPathPattern); err == nil {
			if repository, err := git.PlainClone(clonePath, false, &git.CloneOptions{
				URL:               self.RepositoryURL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			}); err == nil {
				if self.context.dirs == nil {
					self.context.dirs = make(map[string]string)
				}
				self.context.dirs[key] = clonePath

				self.clonePath = clonePath
				return repository, nil
			} else {
				return nil, os.RemoveAll(clonePath)
			}
		} else {
			return nil, err
		}
	}
}

// Utils

func openAndPullGitRepository(path string, repositoryUrl string) (*git.Repository, error) {
	if repository, err := git.PlainOpen(path); err == nil {
		// Pull
		if workTree, err := repository.Worktree(); err == nil {
			if err := workTree.Pull(&git.PullOptions{RemoteName: "origin"}); err != git.NoErrAlreadyUpToDate {
				log.Warningf("could not pull git repository %q from %q: %s", path, repositoryUrl, err.Error())
			}
		} else {
			log.Warningf("could not open git repository %q from %q: %s", path, repositoryUrl, err.Error())
		}

		return repository, nil
	} else {
		return nil, err
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
