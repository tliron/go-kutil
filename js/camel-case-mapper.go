package js

import (
	"reflect"
	"unicode"
)

var CamelCaseMapper camelCaseMapper

type camelCaseMapper struct{}

// goja.FieldNameMapper interface
func (self camelCaseMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return ToCamelCase(f.Name)
}

// goja.FieldNameMapper interface
func (self camelCaseMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return ToCamelCase(m.Name)
}

func ToCamelCase(name string) string {
	runes := []rune(name)
	length := len(runes)
	if (length > 0) && unicode.IsUpper(runes[0]) {
		if (length > 1) && unicode.IsUpper(runes[1]) {
			// If the second rune is also uppercase we'll keep the name as is
			return name
		}
		r := make([]rune, 1, length-1)
		r[0] = unicode.ToLower(runes[0])
		return string(append(r, runes[1:]...))
	}
	return name
}
