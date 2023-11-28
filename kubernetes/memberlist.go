//go:build !wasm

package kubernetes

import (
	"github.com/hashicorp/memberlist"
	"github.com/tliron/commonlog"
	core "k8s.io/api/core/v1"
)

//
// MemberlistPodDiscovery
//

type MemberlistPodDiscovery struct {
	cluster      *memberlist.Memberlist
	podDiscovery *PodDiscovery
	log          commonlog.Logger
}

func StartMemberlistPodDiscovery(cluster *memberlist.Memberlist, namespace string, selector string, frequency float64, log commonlog.Logger) (*MemberlistPodDiscovery, error) {
	self := MemberlistPodDiscovery{
		cluster: cluster,
		log:     log,
	}

	var err error
	if self.podDiscovery, err = StartPodDiscovery(namespace, selector, frequency, self.podsDiscovered, log); err == nil {
		return &self, nil
	} else {
		return nil, err
	}
}

func (self *MemberlistPodDiscovery) Stop() {
	self.podDiscovery.Stop()
}

// kubernetes.PodsDiscoveredFunc signature
func (self *MemberlistPodDiscovery) podsDiscovered(pods []*core.Pod) {
	members := self.cluster.Members()

	var newNodes []string
	for _, pod := range pods {
		found := false
		for _, member := range members {
			if pod.Status.PodIP == member.Addr.String() {
				found = true
				break
			}
		}

		if !found {
			newNodes = append(newNodes, pod.Status.PodIP)
		}
	}

	if len(newNodes) > 0 {
		if self.log != nil {
			self.log.Infof("new pods discovered: %v", newNodes)
		}
		if _, err := self.cluster.Join(newNodes); err != nil {
			if self.log != nil {
				self.log.Errorf("%s", err.Error())
			}
		}
	}
}
