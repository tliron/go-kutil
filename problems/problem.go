package problems

import (
	"fmt"
	"runtime"
	"strings"
)

//
// Caller
//

type Caller struct {
	File     string `json:"file" yaml:"file"`
	Line     int    `json:"line" yaml:"line"`
	Function string `json:"function" yaml:"function"`
}

//
// Problem
//

type Problem struct {
	Section string   `json:"section" yaml:"section"`
	Item    string   `json:"item" yaml:"item"`
	Message string   `json:"message" yaml:"message"`
	Row     int      `json:"row" yaml:"row"`
	Column  int      `json:"column" yaml:"column"`
	Callers []Caller `json:"callers" yaml:"callers"`
}

func NewProblem(section string, item string, message string, row int, column int, skip int) *Problem {
	self := Problem{
		Section: section,
		Item:    item,
		Message: message,
		Row:     row,
		Column:  column,
	}

	callers := make([]uintptr, 1000)
	if count := runtime.Callers(skip+1, callers); count > 0 {
		frames := runtime.CallersFrames(callers)
		for {
			frame, more := frames.Next()

			self.Callers = append(self.Callers, Caller{
				File:     frame.File,
				Line:     frame.Line,
				Function: frame.Function,
			})

			if !more {
				break
			}
		}
	}

	/*
		if _, file, line, ok := runtime.Caller(skip + 1); ok {
			self.Frames = append(self.Frames, Frame{
				SourceFile: file,
				SourceLine: line,
			})
		}
	*/

	return &self
}

// ([fmt.Stringer] interface)
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
