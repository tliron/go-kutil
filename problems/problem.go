package problems

import (
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/terminal"
)

//
// Problem
//

type Problem struct {
	Section    string `json:"section" yaml:"section"`
	Item       string `json:"item" yaml:"item"`
	Message    string `json:"message" yaml:"message"`
	Row        int    `json:"row" yaml:"row"`
	Column     int    `json:"column" yaml:"column"`
	SourceFile string `json:"sourceFile" yaml:"sourceFile"`
	SourceLine int    `json:"sourceLine" yaml:"sourceLine"`
}

func NewProblem(section string, item string, message string, row int, column int, skip int) *Problem {
	self := Problem{
		Section: section,
		Item:    item,
		Message: message,
		Row:     row,
		Column:  column,
	}

	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		self.SourceFile = file
		self.SourceLine = line
	}

	return &self
}

// fmt.Stringer interface
func (self *Problem) String() string {
	r := ""
	if self.Row != -1 {
		r = fmt.Sprintf("@%d", self.Row)
		if self.Column != -1 {
			r += fmt.Sprintf(",%d", self.Column)
		}
		r += " "
	}
	if self.Item != "" {
		r += fmt.Sprintf("%s: ", self.Item)
	}
	r += strings.ReplaceAll(self.Message, "\n", "¶")
	return r
}

func (self *Problem) Equals(problem *Problem) bool {
	return (self.Section == problem.Section) && (self.Item == problem.Item) && (self.Message == problem.Message) && (self.Row == problem.Row) && (self.Column == problem.Column)
}

//
// ProblemSlice
//

type ProblemSlice []*Problem

// sort.Interface

func (self ProblemSlice) Len() int {
	return len(self)
}

func (self ProblemSlice) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

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

//
// Problems
//

type Problems struct {
	Problems ProblemSlice `json:"problems" yaml:"problems"`

	lock sync.RWMutex `json:"-" yaml:"-"`
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
	self.Write(&writer, false, locate)
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

func (self *Problems) ARD() (ard.Value, error) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return ard.Canonicalize(self)
}

func (self *Problems) Write(writer io.Writer, pretty bool, locate bool) bool {
	problems := self.Slice()
	length := len(problems)
	if length > 0 {
		// Sort
		sort.Sort(problems)

		if pretty {
			fmt.Fprintf(writer, "%s (%d)\n", terminal.StyleHeading("Problems"), length)
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
						fmt.Fprintf(writer, "%s\n", terminal.StyleValue(currentSection))
					} else {
						fmt.Fprintf(writer, "%s\n", currentSection)
					}
				} else {
					fmt.Fprintf(writer, "General\n")
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

// Print

func (self *Problems) Print(locate bool) bool {
	return self.Write(terminal.Stderr, true, locate)
}
