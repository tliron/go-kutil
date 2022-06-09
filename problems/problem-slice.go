package problems

import (
	"strings"
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
