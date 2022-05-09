package logging

import "github.com/tliron/kutil/util"

//
// Message
//

type Message interface {
	Set(key string, value any) Message
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

func (self *UnstructuredMessage) Set(key string, value any) Message {
	switch key {
	case "message":
		self.message = util.ToString(value)

	case "scope":
		self.prefix = "{" + util.ToString(value) + "}"

	default:
		if len(self.suffix) > 0 {
			self.suffix += ", "
		}
		self.suffix += key + "=" + util.ToString(value)
	}

	return self
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
