package reflection

import (
	"reflect"

	"github.com/tliron/kutil/util"
)

type EntityTraverser func(any) bool

// Ignore fields tagged with "traverse:ignore" or "lookup"
func TraverseEntities(entity any, locking bool, traverse EntityTraverser) {
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

	if !IsPointerToStruct(reflect.TypeOf(entity)) {
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

		field := value.FieldByName(structField.Name)

		// Ignore unexported fields
		if !field.CanInterface() {
			continue
		}

		fieldType := field.Type()

		if IsPointerToStruct(fieldType) {
			// Compatible with *struct{}
			if !field.IsNil() {
				value := field.Interface()
				if lock != nil {
					lock.Unlock()
				}
				TraverseEntities(value, locking, traverse)
				if lock != nil {
					lock.Lock()
				}
			}
		} else if IsSliceOfPointerToStruct(fieldType) {
			// Compatible with []*struct{}
			if !field.IsNil() {
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
			}
		} else if IsMapOfStringToPointerToStruct(fieldType) {
			// Compatible with map[string]*struct{}
			if !field.IsNil() {
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
	}

	if lock != nil {
		lock.Unlock()
	}
}

//
// EntityWork
//

type EntityWork map[any]struct{}

func (self EntityWork) Start(entityPtr any) bool {
	if _, ok := self[entityPtr]; ok {
		return false
	} else {
		self[entityPtr] = struct{}{}
		return true
	}
}

func (self EntityWork) TraverseEntities(entityPtr any, traverse EntityTraverser) {
	TraverseEntities(entityPtr, false, func(entityPtr any) bool {
		if self.Start(entityPtr) {
			return traverse(entityPtr)
		} else {
			return false
		}
	})
}
