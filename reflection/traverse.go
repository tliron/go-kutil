package reflection

import (
	"reflect"

	"github.com/tliron/kutil/util"
)

type Traverser func(interface{}) bool

// Ignore fields tagged with "traverse:ignore" or "lookup"
func Traverse(object interface{}, traverse Traverser) {
	lock := util.GetLock(object)
	lock.Lock()

	if !traverse(object) {
		lock.Unlock()
		return
	}

	if !IsPtrToStruct(reflect.TypeOf(object)) {
		lock.Unlock()
		return
	}

	value := reflect.ValueOf(object).Elem()

	for _, structField := range GetStructFields(value.Type()) {
		// Has traverse:"ignore" tag?
		traverseTag, ok := structField.Tag.Lookup("traverse")
		if ok && (traverseTag == "ignore") {
			continue
		}

		// Ignore if has "lookup" tag
		if _, ok = structField.Tag.Lookup("lookup"); ok {
			continue
		}

		field := value.FieldByName(structField.Name)
		if !field.CanInterface() {
			// Ignore unexported fields
			continue
		}

		fieldType := field.Type()
		if IsPtrToStruct(fieldType) && !field.IsNil() {
			// Compatible with *struct{}
			lock.Unlock()
			Traverse(field.Interface(), traverse)
			lock.Lock()
		} else if IsSliceOfPtrToStruct(fieldType) {
			// Compatible with []*struct{}
			length := field.Len()
			for index := 0; index < length; index++ {
				element := field.Index(index)
				lock.Unlock()
				Traverse(element.Interface(), traverse)
				lock.Lock()
			}
		} else if IsMapOfStringToPtrToStruct(fieldType) {
			// Compatible with map[string]*struct{}
			for _, mapKey := range field.MapKeys() {
				element := field.MapIndex(mapKey)
				lock.Unlock()
				Traverse(element.Interface(), traverse)
				lock.Lock()
			}
		}
	}

	lock.Unlock()
}
