package sink

import (
	"github.com/hashicorp/memberlist"
	"github.com/tliron/kutil/logging"
)

//
// MemberlistLogEvents
//

type MemberlistLogEvents struct {
	log logging.Logger
}

func NewMemberlistLogEvents(log logging.Logger) *MemberlistLogEvents {
	return &MemberlistLogEvents{log}
}

// memberlist.EventDelegate interface
func (self *MemberlistLogEvents) NotifyJoin(node *memberlist.Node) {
	self.log.Infof("node has joined: %s", node.String())
}

// memberlist.EventDelegate interface
func (self *MemberlistLogEvents) NotifyLeave(node *memberlist.Node) {
	self.log.Infof("node has left: %s", node.String())
}

// memberlist.EventDelegate interface
func (self *MemberlistLogEvents) NotifyUpdate(node *memberlist.Node) {
	self.log.Infof("node was updated: %s", node.String())
}
