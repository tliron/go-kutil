package reflection

import (
	"fmt"
	"reflect"
	"sync"
)

var structFieldsCache sync.Map

// Includes fields "inherited" from anonymous struct fields
// The order of field definition is important! Later fields will override previous fields
func GetStructFields(type_ reflect.Type) []reflect.StructField {
	if structFields, ok := structFieldsCache.Load(type_); ok {
		return structFields.([]reflect.StructField)
	}

	var structFields []reflect.StructField
	length := type_.NumField()
	for index := 0; index < length; index++ {
		structField := type_.Field(index)
		if structField.Anonymous {
			embedded := structField.Type
			for embedded.Kind() == reflect.Pointer {
				embedded = embedded.Elem()
			}
			for _, structField = range GetStructFields(embedded) {
				structFields = appendStructField(structFields, structField)
			}
		} else if structField.IsExported() {
			structFields = appendStructField(structFields, structField)
		}
	}

	structFieldsCache.Store(type_, structFields)

	return structFields
}

func appendStructField(structFields []reflect.StructField, structField reflect.StructField) []reflect.StructField {
	for index, f := range structFields {
		if f.Name == structField.Name {
			// Override
			structFields[index] = structField
			return structFields
		}
	}
	structFields = append(structFields, structField)
	return structFields
}

func GetReferredField(entity reflect.Value, referenceFieldName string, referredFieldName string) (reflect.Value, reflect.Value, bool) {
	referenceField := entity.FieldByName(referenceFieldName)
	if !referenceField.IsValid() {
		panic(fmt.Sprintf("tag refers to unknown field %q in struct: %s", referenceFieldName, entity.Type()))
	}
	if referenceField.Type().Kind() != reflect.Pointer {
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
