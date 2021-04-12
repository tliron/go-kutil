package problems

import (
	"fmt"
	"runtime"
	"strings"
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
