package kubernetes

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
	"github.com/tliron/kutil/util"
	core "k8s.io/api/core/v1"
)

//
// RegistryCredentials
//

type RegistryCredentials struct {
	Username string
	Password string
}

func NewRegistryCredentials(username string, password string) *RegistryCredentials {
	return &RegistryCredentials{
		Username: username,
		Password: password,
	}
}

func NewRegistryCredentialsFromARD(value ard.Value) (*RegistryCredentials, error) {
	node := ard.NewNode(value)
	if v, ok := node.Get("username").String(false); ok {
		var self RegistryCredentials
		self.Username = v
		if v, ok := node.Get("password").String(false); ok {
			self.Password = v
			return &self, nil
		}
	}
	return nil, fmt.Errorf("malformed registry credentials: %s", value)
}

func (self *RegistryCredentials) ToARD() ard.Value {
	return ard.StringMap{
		"username": self.Username,
		"password": self.Password,
	}
}

//
// RegistryCredentialsTable
//

type RegistryCredentialsTable map[string]*RegistryCredentials

func NewRegistryCredentialsTableFromARD(value ard.Value) (RegistryCredentialsTable, error) {
	if map_, ok := value.(ard.Map); ok {
		self := make(RegistryCredentialsTable)
		for server, credentials := range map_ {
			if server_, ok := server.(string); ok {
				var err error
				if self[server_], err = NewRegistryCredentialsFromARD(credentials); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("malformed registry credentials table: %s", value)
			}
		}
		return self, nil
	} else {
		return nil, fmt.Errorf("malformed registry credentials table: %s", value)
	}
}

func NewRegistryCredentialsTableFromSecret(secret *core.Secret) (RegistryCredentialsTable, error) {
	switch secret.Type {
	case core.SecretTypeDockerConfigJson:
		if data, ok := secret.Data[core.DockerConfigJsonKey]; ok {
			if value, err := format.DecodeJSON(util.BytesToString(data)); err == nil {
				if auths := ard.NewNode(value).Get("auths").Data; auths != nil {
					return NewRegistryCredentialsTableFromARD(auths)
				} else {
					return nil, fmt.Errorf("malformed %q secret: %s", core.SecretTypeDockerConfigJson, value)
				}
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("malformed %q secret: %s", core.SecretTypeDockerConfigJson, secret.Data)
		}

	case core.SecretTypeDockercfg:
		if data, ok := secret.Data[core.DockerConfigKey]; ok {
			if value, err := format.DecodeJSON(util.BytesToString(data)); err == nil {
				return NewRegistryCredentialsTableFromARD(value)
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("malformed %q secret: %s", core.SecretTypeDockercfg, secret.Data)
		}
	}

	return nil, nil
}

func (self RegistryCredentialsTable) ToARD() ard.Value {
	servers := make(ard.StringMap)
	for server, credentials := range self {
		servers[server] = credentials.ToARD()
	}
	return servers
}

func (self RegistryCredentialsTable) ToDockerConfigJSON() (string, error) {
	return format.EncodeJSON(ard.StringMap{
		"auths": self.ToARD(),
	}, "")
}

func (self RegistryCredentialsTable) ToSecret(secret *core.Secret) error {
	if dockerConfigJson, err := self.ToDockerConfigJSON(); err == nil {
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

// Utils

func RegistryCredentialsAuth(username string, password string) string {
	// See: https://github.com/kubernetes/kubectl/blob/3874cf79897cfe1e070e592391792658c44b78d4/pkg/generate/versioned/secret_for_docker_registry.go#L166
	auth := fmt.Sprintf("%s:%s", username, password)
	auth = util.ToBase64(util.StringToBytes(auth))
	return auth
}
