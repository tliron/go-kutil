package problems

import (
	"fmt"
)

func (self *Problems) ReportInSection(skip int, message string, section string, row int, column int) bool {
	return self.Append(NewProblem(message, section, row, column, skip+1))
}

func (self *Problems) Report(skip int, message string) bool {
	return self.ReportInSection(skip+1, message, "", -1, -1)
}

func (self *Problems) Reportf(skip int, format string, arg ...interface{}) bool {
	return self.Report(skip+1, fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportProblematic(skip int, problematic Problematic) bool {
	message, section, row, column := problematic.Problem()
	return self.ReportInSection(skip+1, message, section, row, column)
}

func (self *Problems) ReportError(err error) bool {
	if problematic, ok := err.(Problematic); ok {
		return self.ReportProblematic(1, problematic)
	} else {
		return self.Report(1, err.Error())
	}
}
