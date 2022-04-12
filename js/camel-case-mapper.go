package js

import (
	"reflect"

	"github.com/tliron/kutil/util"
)

var CamelCaseMapper camelCaseMapper

type camelCaseMapper struct{}

// goja.FieldNameMapper interface
func (self camelCaseMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return util.ToCamelCase(f.Name)
}

// goja.FieldNameMapper interface
func (self camelCaseMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return util.ToCamelCase(m.Name)
}
