package kubernetes

import (
	contextpkg "context"
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

func GetInternalRegistryHost(context contextpkg.Context, kubernetes kubernetespkg.Interface) (string, error) {
	// TODO: OpenShift

	// Minikube
	if service, err := kubernetes.CoreV1().Services("kube-system").Get(context, "registry", meta.GetOptions{}); err == nil {
		return fmt.Sprintf("%s:80", service.Spec.ClusterIP), nil
	} else {
		return "", err
	}
}

/*
func PushTarballToRegistry(path string, name string) error {
	if tag, err := namepkg.NewTag(name); err == nil {
		if image, err := tarball.ImageFromPath(path, &tag); err == nil {
			return remote.Write(tag, image)
		} else {
			return err
		}
	} else {
		return err
	}
}

func PushGzippedTarballToRegistry(path string, name string) error {
	if tag, err := namepkg.NewTag(name); err == nil {
		opener := func() (io.ReadCloser, error) {
			if reader, err := os.Open(path); err == nil {
				return gzip.NewReader(reader)
			} else {
				return nil, err
			}
		}

		if image, err := tarball.Image(opener, &tag); err == nil {
			return remote.Write(tag, image)
		} else {
			return err
		}
	} else {
		return err
	}
}
*/
