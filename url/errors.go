package url

import "fmt"

//
// NotFound
//

type NotFound struct {
	Message string
}

func NewNotFound(message string) *NotFound {
	return &NotFound{message}
}

func NewNotFoundf(format string, arg ...any) *NotFound {
	return NewNotFound(fmt.Sprintf(format, arg...))
}

// error interface
func (self *NotFound) Error() string {
	return self.Message
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFound)
	return ok
}
