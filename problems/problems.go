package problems

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

//
// ProblemSlice
//

type ProblemSlice []*Problem

// sort.Interface interface
func (self ProblemSlice) Len() int {
	return len(self)
}

// sort.Interface interface
func (self ProblemSlice) Less(i, j int) bool {
	iProblem := self[i]
	jProblem := self[j]
	c := strings.Compare(iProblem.Section, jProblem.Section)
	if c == 0 {
		c = strings.Compare(iProblem.Item, jProblem.Item)
		if c == 0 {
			c = strings.Compare(iProblem.Message, jProblem.Message)
		}
	}
	return c < 0
}

// sort.Interface interface
func (self ProblemSlice) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

//
// Problems
//

type Problems struct {
	Problems ProblemSlice      `json:"problems" yaml:"problems"`
	Stylist  *terminal.Stylist `json:"-" yaml:"-"`

	lock util.RWLocker `json:"-" yaml:"-"`
}

func NewProblems(stylist *terminal.Stylist) *Problems {
	return &Problems{
		Stylist: stylist,
		lock:    util.NewDefaultRWLocker(),
	}
}

func (self *Problems) NewProblems() *Problems {
	return &Problems{
		Stylist: self.Stylist,
		lock:    util.NewDefaultRWLocker(),
	}
}

func (self *Problems) Empty() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return len(self.Problems) == 0
}

func (self *Problems) Append(problem *Problem) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	// Avoid duplicates
	for _, problem_ := range self.Problems {
		if problem.Equals(problem_) {
			return false
		}
	}

	self.Problems = append(self.Problems, problem)
	return true
}

func (self *Problems) Merge(problems *Problems) bool {
	if self == problems {
		// Merging into self
		return false
	}

	merged := false
	problems.lock.RLock()
	defer problems.lock.RUnlock()
	for _, problem := range problems.Problems {
		if self.Append(problem) {
			merged = true
		}
	}

	return merged
}

func (self *Problems) ToString(locate bool) string {
	var writer strings.Builder
	self.Write(&writer, nil, false, locate)
	return strings.TrimRight(writer.String(), "\n")
}

// fmt.Stringer interface
func (self *Problems) String() string {
	return self.ToString(false)
}

func (self *Problems) Slice() ProblemSlice {
	problems := make(ProblemSlice, len(self.Problems))
	self.lock.RLock()
	copy(problems, self.Problems)
	self.lock.RUnlock()
	return problems
}

func (self *Problems) Write(writer io.Writer, stylist *terminal.Stylist, pretty bool, locate bool) bool {
	problems := self.Slice()
	length := len(problems)
	if length > 0 {
		if stylist == nil {
			stylist = terminal.NewStylist(false)
		}

		// Sort
		sort.Sort(problems)

		if pretty {
			fmt.Fprintf(writer, "%s (%d)\n", stylist.Heading("Problems"), length)
		} else {
			fmt.Fprintf(writer, "%s (%d)\n", "Problems", length)
		}

		var currentSection string
		for _, problem := range problems {
			section := problem.Section
			if currentSection != section {
				currentSection = section
				fmt.Fprint(writer, terminal.IndentString(1))
				if currentSection != "" {
					if pretty {
						fmt.Fprintf(writer, "%s\n", stylist.Value(currentSection))
					} else {
						fmt.Fprintf(writer, "%s\n", currentSection)
					}
				} else {
					fmt.Fprintln(writer, "General")
				}
			}

			fmt.Fprint(writer, terminal.IndentString(2))
			fmt.Fprintf(writer, "%s\n", problem)

			if locate && (problem.SourceFile != "") {
				fmt.Fprint(writer, terminal.IndentString(2))
				fmt.Fprintf(writer, "└─%s:%d\n", problem.SourceFile, problem.SourceLine)
			}
		}
		return true
	}
	return false
}

func (self *Problems) WithError(err error, locate bool) error {
	var writer strings.Builder
	if err != nil {
		writer.WriteString(err.Error())
	}
	if len(self.Problems) > 0 {
		if writer.Len() > 0 {
			writer.WriteRune('\n')
		}
		writer.WriteString(self.ToString(locate))
	}
	return errors.New(writer.String())
}

// Print

func (self *Problems) Print(locate bool) bool {
	return self.Write(terminal.Stderr, terminal.Stylize, true, locate)
}
