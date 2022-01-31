package kubernetes

import (
	contextpkg "context"
	"time"

	"github.com/tliron/kutil/logging"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedcore "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

//
// PodDiscovery
//

type PodsDiscoveredFunc func([]*core.Pod)

type PodDiscovery struct {
	selector       string
	context        contextpkg.Context
	podsDiscovered PodsDiscoveredFunc
	log            logging.Logger

	pods    typedcore.PodInterface
	started chan struct{}
}

func StartPodDiscovery(namespace string, selector string, frequency float64, podsDiscovered PodsDiscoveredFunc, log logging.Logger) (*PodDiscovery, error) {
	self := PodDiscovery{
		selector:       selector,
		context:        contextpkg.TODO(),
		podsDiscovered: podsDiscovered,
		log:            log,
	}

	if config, err := rest.InClusterConfig(); err == nil {
		if client, err := kubernetes.NewForConfig(config); err == nil {
			self.pods = client.CoreV1().Pods(namespace)
			self.start(frequency)
			return &self, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *PodDiscovery) Stop() {
	close(self.started)
}

func (self *PodDiscovery) start(frequency float64) {
	ticker := time.NewTicker(time.Duration(frequency * 1000000000.0)) // seconds to nanoseconds
	go func() {
		for {
			select {
			case <-ticker.C:
				self.list()

			case <-self.started:
				ticker.Stop()
				return
			}
		}
	}()
}

func (self *PodDiscovery) list() {
	//self.log.Debugf("listing nodes: %s", self.selector)
	if pods, err := self.pods.List(self.context, meta.ListOptions{LabelSelector: self.selector}); err == nil {
		var list []*core.Pod
		for _, pod := range pods.Items {
			list = append(list, &pod)
		}
		self.podsDiscovered(list)
	} else if !errors.IsNotFound(err) {
		if self.log != nil {
			self.log.Errorf("%s", err.Error())
		}
	}
}
