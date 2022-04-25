package sink

import (
	"log"
	"strings"

	"github.com/hashicorp/memberlist"
	"github.com/tliron/kutil/logging"
)

const (
	MEMBERLIST_ERR_PREFIX   = "[ERR] memberlist: "
	MEMBERLIST_WARN_PREFIX  = "[WARN] memberlist: "
	MEMBERLIST_DEBUG_PREFIX = "[DEBUG] memberlist: "
)

func NewMemberlistStandardLog(name []string) *log.Logger {
	return NewStandardLogger(func(line string) logging.Message {
		level := logging.Debug

		if strings.HasPrefix(line, MEMBERLIST_ERR_PREFIX) {
			line = line[len(MEMBERLIST_ERR_PREFIX):]
			level = logging.Error
		} else if strings.HasPrefix(line, MEMBERLIST_WARN_PREFIX) {
			line = line[len(MEMBERLIST_WARN_PREFIX):]
			level = logging.Warning
		} else if strings.HasPrefix(line, MEMBERLIST_DEBUG_PREFIX) {
			line = line[len(MEMBERLIST_DEBUG_PREFIX):]
		}

		if message := logging.NewMessage(name, level, 2); message != nil {
			message.Set("message", line)
			return message
		} else {
			return nil
		}
	})
}

//
// MemberlistEventLog
//

type MemberlistEventLog struct {
	log logging.Logger
}

func NewMemberlistEventLog(log logging.Logger) *MemberlistEventLog {
	return &MemberlistEventLog{log}
}

// memberlist.EventDelegate interface
func (self *MemberlistEventLog) NotifyJoin(node *memberlist.Node) {
	self.log.Infof("node has joined: %s", node.String())
}

// memberlist.EventDelegate interface
func (self *MemberlistEventLog) NotifyLeave(node *memberlist.Node) {
	self.log.Infof("node has left: %s", node.String())
}

// memberlist.EventDelegate interface
func (self *MemberlistEventLog) NotifyUpdate(node *memberlist.Node) {
	self.log.Infof("node was updated: %s", node.String())
}
