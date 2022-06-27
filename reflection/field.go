package reflection

import (
	"fmt"
	"reflect"
	"sync"
)

var structFieldsCache sync.Map

// Includes fields "inherited" from anonymous struct pointer fields
// The order of field definition is important! Later fields will override previous fields
func GetStructFields(type_ reflect.Type) []reflect.StructField {
	if v, ok := structFieldsCache.Load(type_); ok {
		return v.([]reflect.StructField)
	}

	var structFields []reflect.StructField
	fields := type_.NumField()
	for index := 0; index < fields; index++ {
		structField := type_.Field(index)
		if structField.Anonymous && (structField.Type.Kind() == reflect.Ptr) {
			for _, structField = range GetStructFields(structField.Type.Elem()) {
				structFields = appendStructField(structFields, structField)
			}
		} else {
			structFields = appendStructField(structFields, structField)
		}
	}

	structFieldsCache.Store(type_, structFields)

	return structFields
}

func appendStructField(structFields []reflect.StructField, structField reflect.StructField) []reflect.StructField {
	found := false
	for index, f := range structFields {
		if f.Name == structField.Name {
			// Override
			structFields[index] = structField
			found = true
			break
		}
	}
	if !found {
		structFields = append(structFields, structField)
	}
	return structFields
}

func GetReferredField(entity reflect.Value, referenceFieldName string, referredFieldName string) (reflect.Value, reflect.Value, bool) {
	referenceField := entity.FieldByName(referenceFieldName)
	if !referenceField.IsValid() {
		panic(fmt.Sprintf("tag refers to unknown field %q in struct: %s", referenceFieldName, entity.Type()))
	}
	if referenceField.Type().Kind() != reflect.Ptr {
		panic(fmt.Sprintf("tag refers to non-pointer field %q in struct: %s", referenceFieldName, entity.Type()))
	}

	if referenceField.IsNil() {
		// Reference is empty
		return referenceField, reflect.Value{}, false
	}

	referredField := referenceField.Elem().FieldByName(referredFieldName)

	if !referredField.IsValid() {
		panic(fmt.Sprintf("tag's field name %q not found in the entity referred to by %q in struct: %s", referredFieldName, referenceFieldName, entity.Type()))
	}

	return referenceField, referredField, true
}
