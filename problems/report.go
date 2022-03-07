package problems

import (
	"fmt"
)

func (self *Problems) ReportFull(skip int, section string, item string, message string, row int, column int) bool {
	return self.Append(NewProblem(section, item, message, row, column, skip+1))
}

func (self *Problems) Report(skip int, item string, message string) bool {
	return self.ReportFull(skip+1, "", item, message, -1, -1)
}

func (self *Problems) Reportf(skip int, item string, format string, arg ...any) bool {
	return self.Report(skip+1, item, fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportProblematic(skip int, problematic Problematic) bool {
	section, item, message, row, column := problematic.Problem(self.Stylist)
	return self.ReportFull(skip+1, section, item, message, row, column)
}

func (self *Problems) ReportError(err error) bool {
	if problematic, ok := err.(Problematic); ok {
		return self.ReportProblematic(1, problematic)
	} else {
		return self.Report(1, "", err.Error())
	}
}
