package reflection

import (
	"reflect"
)

func GetTaggedFields(structPtr any, name string) []reflect.Value {
	var fields []reflect.Value
	value := reflect.ValueOf(structPtr).Elem()
	for fieldName := range GetFieldTagsForValue(value, name) {
		field := value.FieldByName(fieldName)
		fields = append(fields, field)
	}
	return fields
}

func GetFieldTagsForValue(value reflect.Value, name string) map[string]string {
	return GetFieldTagsForType(value.Type(), name)
}

func GetFieldTagsForType(type_ reflect.Type, name string) map[string]string {
	tags := make(map[string]string)
	for _, structField := range GetStructFields(type_) {
		if value, ok := structField.Tag.Lookup(name); ok {
			tags[structField.Name] = value
		}
	}
	return tags
}
