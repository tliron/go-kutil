package sink

import (
	"github.com/hashicorp/memberlist"
	"github.com/tliron/kutil/logging"
)

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
