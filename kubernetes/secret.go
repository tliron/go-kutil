package kubernetes

import (
	contextpkg "context"
	"fmt"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

func CopySecret(context contextpkg.Context, kubernetes kubernetespkg.Interface, fromNamespace string, fromName string, toNamespace string, toName string) (*core.Secret, error) {
	if secret, err := kubernetes.CoreV1().Secrets(fromNamespace).Get(context, fromName, meta.GetOptions{}); err == nil {
		secret = secret.DeepCopy()
		secret.ResourceVersion = ""
		secret.Namespace = toNamespace
		secret.Name = toName

		if err := kubernetes.CoreV1().Secrets(secret.Namespace).Delete(context, secret.Name, meta.DeleteOptions{}); err == nil {
			return kubernetes.CoreV1().Secrets(secret.Namespace).Create(context, secret, meta.CreateOptions{})
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func GetSecretTLSCertBytes(secret *core.Secret, secretDataKey string) ([]byte, error) {
	switch secret.Type {
	case core.SecretTypeTLS:
		if secretDataKey == "" {
			secretDataKey = core.TLSCertKey
		}

		if bytes, ok := secret.Data[secretDataKey]; ok {
			return bytes, nil
		} else {
			return nil, fmt.Errorf("no data key %q in %q secret: %s", secretDataKey, secret.Type, secret.Data)
		}

	case core.SecretTypeServiceAccountToken:
		if secretDataKey == "" {
			secretDataKey = core.ServiceAccountRootCAKey
		}

		if bytes, ok := secret.Data[secretDataKey]; ok {
			return bytes, nil
		} else {
			return nil, fmt.Errorf("no data key %q in %q secret: %s", secretDataKey, secret.Type, secret.Data)
		}

	default:
		return nil, fmt.Errorf("unsupported TLS secret type: %s", secret.Type)
	}
}

/*
func GetSecretTLSCertPool(secret *core.Secret, secretDataKey string) (*x509.CertPool, error) {
	switch secret.Type {
	case core.SecretTypeTLS:
		if secretDataKey == "" {
			secretDataKey = core.TLSCertKey
		}

		if bytes, ok := secret.Data[secretDataKey]; ok {
			return util.ParseX509CertPool(bytes)
		} else {
			return nil, fmt.Errorf("no data key %q in %q secret: %s", secretDataKey, secret.Type, secret.Data)
		}

	case core.SecretTypeServiceAccountToken:
		if secretDataKey == "" {
			secretDataKey = core.ServiceAccountRootCAKey
		}

		if bytes, ok := secret.Data[secretDataKey]; ok {
			return util.ParseX509CertPool(bytes)
		} else {
			return nil, fmt.Errorf("no data key %q in %q secret: %s", secretDataKey, secret.Type, secret.Data)
		}

	default:
		return nil, fmt.Errorf("unsupported TLS secret type: %s", secret.Type)
	}
}
*/
