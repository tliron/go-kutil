package reflection

import (
	"reflect"

	"github.com/tliron/kutil/util"
)

type EntityTraverser func(interface{}) bool

// Ignore fields tagged with "traverse:ignore" or "lookup"
func TraverseEntities(entity interface{}, locking bool, traverse EntityTraverser) {
	var lock util.RWLocker
	if locking {
		lock = util.GetEntityLock(entity)
	}
	if lock != nil {
		lock.Lock()
	}

	if !traverse(entity) {
		if lock != nil {
			lock.Unlock()
		}
		return
	}

	if !IsPtrToStruct(reflect.TypeOf(entity)) {
		if lock != nil {
			lock.Unlock()
		}
		return
	}

	value := reflect.ValueOf(entity).Elem()

	for _, structField := range GetStructFields(value.Type()) {
		// Ignore if has traverse:"ignore" tag
		if traverseTag, ok := structField.Tag.Lookup("traverse"); ok && (traverseTag == "ignore") {
			continue
		}

		// Ignore if has "lookup" tag
		// TODO: should this be *here*?
		if _, ok := structField.Tag.Lookup("lookup"); ok {
			continue
		}

		field := value.FieldByName(structField.Name)

		// Ignore unexported fields
		if !field.CanInterface() {
			continue
		}

		fieldType := field.Type()

		if IsPtrToStruct(fieldType) && !field.IsNil() {
			// Compatible with *struct{}
			value := field.Interface()
			if lock != nil {
				lock.Unlock()
			}
			TraverseEntities(value, locking, traverse)
			if lock != nil {
				lock.Lock()
			}
		} else if IsSliceOfPtrToStruct(fieldType) {
			// Compatible with []*struct{}
			length := field.Len()
			elements := make([]reflect.Value, length)
			for index := 0; index < length; index++ {
				elements[index] = field.Index(index)
			}

			for _, element := range elements {
				value := element.Interface()
				if lock != nil {
					lock.Unlock()
				}
				TraverseEntities(value, locking, traverse)
				if lock != nil {
					lock.Lock()
				}
			}
		} else if IsMapOfStringToPtrToStruct(fieldType) {
			// Compatible with map[string]*struct{}
			keys := field.MapKeys()
			elements := make([]reflect.Value, len(keys))
			for index, key := range keys {
				elements[index] = field.MapIndex(key)
			}

			for _, element := range elements {
				value := element.Interface()
				if lock != nil {
					lock.Unlock()
				}
				TraverseEntities(value, locking, traverse)
				if lock != nil {
					lock.Lock()
				}
			}
		}
	}

	if lock != nil {
		lock.Unlock()
	}
}
