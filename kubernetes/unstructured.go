package kubernetes

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	discoverypkg "k8s.io/client-go/discovery"
)

func NewUnstructuredFromYAMLTemplate(code string, data any) (*unstructured.Unstructured, error) {
	if object, err := ard.DecodeYAMLTemplate(code, data); err == nil {
		object, _ = ard.NormalizeStringMaps(object)
		if stringMap, ok := object.(ard.StringMap); ok {
			return &unstructured.Unstructured{Object: stringMap}, nil
		} else {
			return nil, fmt.Errorf("not a map[string]any: %T", object)
		}
	} else {
		return nil, err
	}
}

// controllerObject must also support schema.ObjectKind interface
func SetControllerOfUnstructured(object *unstructured.Unstructured, controllerObject meta.Object) error {
	if controllerObjectKind, ok := controllerObject.(schema.ObjectKind); ok {
		ownerReferences := object.GetOwnerReferences()
		ownerReferences = append(ownerReferences, *meta.NewControllerRef(controllerObject, controllerObjectKind.GroupVersionKind()))
		object.SetOwnerReferences(ownerReferences)
		return nil
	} else {
		return fmt.Errorf("controller object does not support schema.ObjectKind interface: %+v", controllerObject)
	}
}

func GetUnstructuredGVK(object *unstructured.Unstructured) (schema.GroupVersionKind, error) {
	return ParseGVK(object.GetAPIVersion(), object.GetKind())
}

func FindResourceForUnstructured(discovery discoverypkg.DiscoveryInterface, object *unstructured.Unstructured, supportedVerbs ...string) (schema.GroupVersionResource, error) {
	if gvk, err := ParseGVK(object.GetAPIVersion(), object.GetKind()); err == nil {
		return FindResourceForKind(discovery, gvk, supportedVerbs...)
	} else {
		return schema.GroupVersionResource{}, err
	}
}

//
// UnstructuredResourceEventHandler
//

type OnAddedFunc = func(object *unstructured.Unstructured) error

type OnUpdatedFunc = func(oldObject *unstructured.Unstructured, newObject *unstructured.Unstructured) error

type OnDeletedFunc = func(object *unstructured.Unstructured) error

type UnstructuredResourceEventHandler struct {
	onAdded   OnAddedFunc
	onUpdated OnUpdatedFunc
	onDeleted OnDeletedFunc
}

func NewUnstructuredResourceEventHandler(onAdded OnAddedFunc, onUpdated OnUpdatedFunc, onDeleted OnDeletedFunc) *UnstructuredResourceEventHandler {
	return &UnstructuredResourceEventHandler{
		onAdded:   onAdded,
		onUpdated: onUpdated,
		onDeleted: onDeleted,
	}
}

// cache.ResourceEventHandler interface
func (self *UnstructuredResourceEventHandler) OnAdd(object any) {
	utilruntime.HandleError(self.onAdded(object.(*unstructured.Unstructured)))
}

// cache.ResourceEventHandler interface
func (self *UnstructuredResourceEventHandler) OnUpdate(oldObject any, newObject any) {
	utilruntime.HandleError(self.onUpdated(oldObject.(*unstructured.Unstructured), newObject.(*unstructured.Unstructured)))
}

// cache.ResourceEventHandler interface
func (self *UnstructuredResourceEventHandler) OnDelete(object any) {
	utilruntime.HandleError(self.onDeleted(object.(*unstructured.Unstructured)))
}
