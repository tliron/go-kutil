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

func (self *Message) Set(key string, value interface{}) logging.Message {
	switch value_ := value.(type) {
	case string:
		self.event.Str(key, value_)
	case int:
		self.event.Int(key, value_)
	case int64:
		self.event.Int64(key, value_)
	case int32:
		self.event.Int32(key, value_)
	case int16:
		self.event.Int16(key, value_)
	case int8:
		self.event.Int8(key, value_)
	case uint:
		self.event.Uint(key, value_)
	case uint64:
		self.event.Uint64(key, value_)
	case uint32:
		self.event.Uint32(key, value_)
	case uint16:
		self.event.Uint16(key, value_)
	case uint8:
		self.event.Uint8(key, value_)
	case float64:
		self.event.Float64(key, value_)
	case float32:
		self.event.Float32(key, value_)
	case bool:
		self.event.Bool(key, value_)
	case []byte:
		self.event.Bytes(key, value_)
	case fmt.Stringer:
		self.event.Stringer(key, value_)
	default:
		self.event.Interface(key, value_)
	}

	return self
}

func (self *Message) Send() {
	self.event.Send()
}
