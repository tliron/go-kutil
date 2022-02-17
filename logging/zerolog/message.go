package zerolog

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/tliron/kutil/logging"
)

//
// Message
//

type Message struct {
	event *zerolog.Event
}

func NewMessage(event *zerolog.Event) logging.Message {
	return &Message{event: event}
}

// logging.Message interface

func (self *Message) Set(name string, value interface{}) {
	switch value_ := value.(type) {
	case string:
		self.event.Str(name, value_)
	case int:
		self.event.Int(name, value_)
	case int64:
		self.event.Int64(name, value_)
	case int32:
		self.event.Int32(name, value_)
	case int16:
		self.event.Int16(name, value_)
	case int8:
		self.event.Int8(name, value_)
	case uint:
		self.event.Uint(name, value_)
	case uint64:
		self.event.Uint64(name, value_)
	case uint32:
		self.event.Uint32(name, value_)
	case uint16:
		self.event.Uint16(name, value_)
	case uint8:
		self.event.Uint8(name, value_)
	case float64:
		self.event.Float64(name, value_)
	case float32:
		self.event.Float32(name, value_)
	case bool:
		self.event.Bool(name, value_)
	case []byte:
		self.event.Bytes(name, value_)
	case fmt.Stringer:
		self.event.Stringer(name, value_)
	default:
		self.event.Interface(name, value_)
	}
}

func (self *Message) Send() {
	self.event.Send()
}
