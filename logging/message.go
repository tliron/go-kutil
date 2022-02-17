package logging

import "fmt"

//
// Message
//

type Message interface {
	Set(name string, value interface{})
	Send()
}

//
// UnstructuredMessage
//

type SendUnstructuredMessageFunc func(message string)

type UnstructuredMessage struct {
	prefix  string
	message string
	suffix  string
	send    SendUnstructuredMessageFunc
}

func NewUnstructuredMessage(send SendUnstructuredMessageFunc) *UnstructuredMessage {
	return &UnstructuredMessage{
		send: send,
	}
}

// Message interface

func (self *UnstructuredMessage) Set(name string, value interface{}) {
	switch name {
	case "domain":
		self.prefix = "{" + toString(value) + "}"

	case "message":
		self.message = toString(value)

	default:
		if len(self.suffix) > 0 {
			self.suffix += ", "
		}
		self.suffix += name + "=" + toString(value)
	}
}

func (self *UnstructuredMessage) Send() {
	message := self.prefix

	if len(self.message) > 0 {
		if len(message) > 0 {
			message += " "
		}
		message += self.message
	}

	if len(self.suffix) > 0 {
		if len(message) > 0 {
			message += " "
		}
		message += self.suffix
	}

	self.send(message)
}

func toString(value interface{}) string {
	switch value_ := value.(type) {
	case string:
		return value_
	case fmt.Stringer:
		return value_.String()
	default:
		return fmt.Sprintf("%s", value_)
	}
}
