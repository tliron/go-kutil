package kubernetes

import (
	contextpkg "context"
	"time"

	"github.com/tliron/commonlog"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	errorspkg "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

var interval = 10 * time.Second
var timeout = 60 * time.Second

func WaitForPod(context contextpkg.Context, kubernetes kubernetespkg.Interface, log commonlog.Logger, namespace string, appName string) (*core.Pod, error) {
	log.Infof("waiting for a pod for app %q", appName)

	var pod *core.Pod

	err := Wait(context, func(context contextpkg.Context) (bool, error) {
		if pods, err := GetPods(context, kubernetes, namespace, appName); err == nil {
			for _, pod_ := range pods.Items {
				for _, containerStatus := range pod_.Status.ContainerStatuses {
					if containerStatus.Ready {
						log.Infof("container %q ready for pod %q", containerStatus.Name, pod_.Name)
					} else {
						return false, nil
					}
				}

				for _, condition := range pod_.Status.Conditions {
					switch condition.Type {
					case core.ContainersReady:
						if condition.Status == core.ConditionTrue {
							pod = &pod_
							return true, nil
						}
					}
				}
			}
			return false, nil
		} else if errorspkg.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	})

	if (err == nil) && (pod != nil) {
		log.Infof("a pod is available for app %q", appName)
		return pod, nil
	} else {
		return nil, err
	}
}

func WaitForDeployment(context contextpkg.Context, kubernetes kubernetespkg.Interface, log commonlog.Logger, namespace string, appName string) (*apps.Deployment, error) {
	log.Infof("waiting for a deployment for app %q", appName)

	var deployment *apps.Deployment
	err := Wait(context, func(context contextpkg.Context) (bool, error) {
		var err error
		if deployment, err = kubernetes.AppsV1().Deployments(namespace).Get(context, appName, meta.GetOptions{}); err == nil {
			for _, condition := range deployment.Status.Conditions {
				switch condition.Type {
				case apps.DeploymentAvailable:
					if condition.Status == core.ConditionTrue {
						return true, nil
					}

				case apps.DeploymentReplicaFailure:
					if condition.Status == core.ConditionTrue {
						log.Infof("replica failure for a deployment for app %q", appName)
					}
				}
			}
			return false, nil
		} else {
			return false, err
		}
	})

	if err == nil {
		log.Infof("a deployment is available for app %q", appName)
		return deployment, nil
	} else {
		return nil, err
	}
}

func WaitForDeletion(context contextpkg.Context, log commonlog.Logger, name string, condition func() bool) {
	if err := Wait(context, func(context contextpkg.Context) (bool, error) {
		log.Infof("waiting for %s to delete", name)
		return !condition(), nil
	}); err != nil {
		log.Warningf("error while waiting for %s to delete: %s", name, err.Error())
	}
}

func Wait(context contextpkg.Context, condition wait.ConditionWithContextFunc) error {
	return wait.PollUntilContextTimeout(context, interval, timeout, true, condition)
}

/*
func (self *Client) waitForPodContainers(appName string, deployment *apps.Deployment) error {
	log.Infof("waiting for pods for %q", appName)

	return waitpkg.PollImmediate(time.Second, timeout, func() (bool, error) {
		if pods, err := kubernetes.GetPods(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
			for _, pod := range pods.Items {
				if self.isPodOwnedBy(&pod, deployment) {
					for _, container := range pod.Spec.Containers {
						if err := self.Exec(self.Namespace, pod.Name, container.Name, nil, nil, "echo"); err == nil {
							log.Infof("container %q available for pod: %s", container.Name, pod.Name)
						} else {
							return false, nil
						}
					}
				}
				log.Infof("pod available for %q: %s", appName, pod.Name)
				return true, nil
			}
			return false, nil
		} else {
			return false, err
		}
	})
}

func (self *Client) waitForAPod(appName string, deployment *apps.Deployment) error {
	log.Infof("waiting for a pod for %q", appName)

	return waitpkg.PollImmediate(time.Second, timeout, func() (bool, error) {
		if pods, err := kubernetes.GetPods(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
			for _, pod := range pods.Items {
				if self.isPodOwnedBy(&pod, deployment) {
					for _, containerStatus := range pod.Status.ContainerStatuses {
						if containerStatus.Ready {
							log.Infof("container %q ready for pod: %s", containerStatus.Name, pod.Name)
						} else {
							return false, nil
						}
					}

					for _, condition := range pod.Status.Conditions {
						switch condition.Type {
						case core.ContainersReady:
							if condition.Status == core.ConditionTrue {
								log.Infof("pod ready for %q: %s", appName, pod.Name)
								return true, nil
							}
						}
					}
				}
			}
			return false, nil
		} else {
			return false, err
		}
	})
}

func (self *Client) isPodOwnedBy(pod *core.Pod, deployment *apps.Deployment) bool {
	for _, owner := range pod.OwnerReferences {
		if (owner.APIVersion == "apps/v1") && (owner.Kind == "ReplicaSet") {
			if replicaSet, err := self.Kubernetes.AppsV1().ReplicaSets(self.Namespace).Get(self.Context, owner.Name, meta.GetOptions{}); err == nil {
				if self.isReplicaSetOwnedBy(replicaSet, deployment) {
					return true
				}
			}
		}
	}
	return false
}

func (self *Client) isReplicaSetOwnedBy(replicaSet *apps.ReplicaSet, deployment *apps.Deployment) bool {
	for _, owner := range replicaSet.OwnerReferences {
		if owner.UID == deployment.UID {
			return true
		}
	}
	return false
}
*/
