package kubernetes

import (
	"fmt"

	"github.com/tliron/kutil/format"
	"github.com/tliron/kutil/util"
	core "k8s.io/api/core/v1"
)

func SetSecretDockerConfigJson(secret *core.Secret, server string, username string, password string) error {
	if dockerConfigJson, err := DockerConfigJson(server, username, password); err == nil {
		secret.Type = core.SecretTypeDockerConfigJson
		secret.Data = map[string][]byte{
			// TODO: do we need to call ToBase64?
			core.DockerConfigJsonKey: util.StringToBytes(dockerConfigJson),
		}
		return nil
	} else {
		return err
	}
}

func DockerConfig(server string, username string, password string) interface{} {
	return map[string]map[string]map[string]string{
		"auths": map[string]map[string]string{
			server: map[string]string{
				"username": username,
				"password": password,
				"auth":     DockerConfigAuth(username, password),
			},
		},
	}
}

func DockerConfigAuth(username string, password string) string {
	// See: https://github.com/kubernetes/kubectl/blob/3874cf79897cfe1e070e592391792658c44b78d4/pkg/generate/versioned/secret_for_docker_registry.go#L166
	auth := fmt.Sprintf("%s:%s", username, password)
	auth = util.ToBase64(util.StringToBytes(auth))
	return auth
}

func DockerConfigJson(server string, username string, password string) (string, error) {
	return format.EncodeJSON(DockerConfig(server, username, password), "")
}
