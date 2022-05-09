package journal

import (
	"strings"

	"github.com/coreos/go-systemd/journal"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/util"
)

//
// Message
//

type Message struct {
	priority journal.Priority
	prefix   string
	message  string
	vars     map[string]string
}

func NewMessage(priority journal.Priority, prefix string) logging.Message {
	return &Message{
		priority: priority,
		prefix:   prefix,
	}
}

// logging.Message interface
func (self *Message) Set(key string, value any) logging.Message {
	switch key {
	case "message":
		self.message = util.ToString(value)

	default:
		// See: https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html
		key = strings.ToUpper(key)
		if self.vars == nil {
			self.vars = make(map[string]string)
		}
		self.vars[key] = util.ToString(value)
	}

	return self
}

// logging.Message interface
func (self *Message) Send() {
	journal.Send(self.prefix+self.message, self.priority, self.vars)
}
