package kubernetes

import (
	"fmt"
	"net/http"

	errorspkg "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Variation on errorspkg.NewNotFound
func NewNotFound(gk schema.GroupKind, message string) *errorspkg.StatusError {
	return &errorspkg.StatusError{meta.Status{
		Status: meta.StatusFailure,
		Code:   http.StatusNotFound,
		Reason: meta.StatusReasonNotFound,
		Details: &meta.StatusDetails{
			Group: gk.Group,
			Kind:  gk.Kind,
		},
		Message: fmt.Sprintf("%s not found: %s", gk.String(), message),
	}}
}
