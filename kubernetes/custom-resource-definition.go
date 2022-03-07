package kubernetes

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func JSONString(value any) apiextensions.JSON {
	return apiextensions.JSON{
		Raw: []byte(fmt.Sprintf("%q", value)),
	}
}
